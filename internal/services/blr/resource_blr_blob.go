package blr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/blr"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceBLRBlob() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceBLRBlobCreate,
		ReadContext:   resourceBLRBlobRead,
		UpdateContext: resourceBLRBlobUpdate,
		DeleteContext: resourceBLRBlobDelete,

		Schema: map[string]*schema.Schema{
			"data_type_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"blob_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"blob_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_access_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_access_url_expiry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"policy": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     policySchema(),
			},
		},
	}
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func policySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"statement": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 10,
				Elem:     statementSchema(),
			},
		},
	}
}

func principalSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hsdp": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 10,
				Elem:     tools.StringSchema(),
			},
		},
	}
}

func statementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"principal": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem:     principalSchema(),
			},
			"effect": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeSet,
				MaxItems: 32,
				Required: true,
				Elem:     tools.StringSchema(),
			},
		},
	}
}

func schemaToBlob(d *schema.ResourceData) blr.Blob {
	dataTypeName := d.Get("data_type_name").(string)
	blobPath := d.Get("blob_path").(string)
	blobName := d.Get("blob_name").(string)
	virtualPath := d.Get("virtual_path").(string)
	virtualName := d.Get("virtual_name").(string)
	tagList := d.Get("tags").(map[string]interface{})
	tags := make([]blr.Tag, 0)
	for t, v := range tagList {
		if val, ok := v.(string); ok {
			tags = append(tags, blr.Tag{
				Key:   t,
				Value: val,
			})
		}
	}

	resource := blr.Blob{
		DataType:    dataTypeName,
		BlobName:    blobName,
		BlobPath:    blobPath,
		VirtualPath: virtualPath,
		VirtualName: virtualName,
	}
	if len(tags) > 0 {
		resource.Tags = &tags
	}
	return resource
}

func schemaToPolicy(d *schema.ResourceData) blr.BlobPolicy {
	resource := blr.BlobPolicy{
		ResourceType: "AccessPolicy",
	}

	var statements []blr.PolicyStatement

	if v, ok := d.GetOk("policy"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL { // Statements
			mVi := vi.(map[string]interface{})
			statementsList := mVi["statement"].(*schema.Set)
			sL := statementsList.List()
			for _, si := range sL {
				mSi := si.(map[string]interface{})
				var statement blr.PolicyStatement
				statement.ResourceType = "AccessPolicy"
				statement.Effect = mSi["effect"].(string)
				statement.Action = tools.ExpandStringList(mSi["action"].(*schema.Set).List())
				principlesList := mSi["principal"].(*schema.Set).List()
				for _, pl := range principlesList {
					pr := pl.(map[string]interface{})
					statement.Principal.HSDP = tools.ExpandStringList(pr["hsdp"].(*schema.Set).List())
				}
				statements = append(statements, statement)
			}

		}
		resource.Statement = statements
	}
	return resource
}

func blobToSchema(resource blr.Blob, d *schema.ResourceData) {
	_ = d.Set("data_type_name", resource.DataType)
}

func policyToSchema(resource blr.PolicyStatement, d *schema.ResourceData) {

}

func resourceBLRBlobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.BLRClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToBlob(d)
	policy := schemaToPolicy(d)

	var created *blr.Blob
	var resp *blr.Response
	err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Blobs.Create(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if created == nil {
		return diag.FromErr(fmt.Errorf("failed to create resource: %d", resp.StatusCode))
	}
	if len(policy.Statement) > 0 {
		err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
			var err error
			_, resp, err = client.Blobs.SetPolicy(*created, policy)
			if err != nil {
				_ = client.TokenRefresh()
			}
			if resp == nil {
				return nil, err
			}
			return resp.Response, err
		})
	}
	if err != nil {
		_, _, _ = client.Blobs.Delete(*created)
		return diag.FromErr(err)
	}
	d.SetId(created.ID)
	return resourceBLRBlobRead(ctx, d, m)
}

func resourceBLRBlobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.BLRClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	var resource *blr.Blob
	var resp *blr.Response
	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Blobs.GetByID(id)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	blobToSchema(*resource, d)
	return diags
}

func resourceBLRBlobUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	_, err := c.BLRClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()

	resource := schemaToBlob(d)
	resource.ID = id

	// TODO: figure out what to do here
	return diags
}

func resourceBLRBlobDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.BLRClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	resource, _, err := client.Blobs.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.Blobs.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
