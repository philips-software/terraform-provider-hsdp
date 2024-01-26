package blr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/connect/blr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func importStatePassthroughSetGuidContext(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	var id string
	count, _ := fmt.Sscanf(d.Id(), "BlobStorePolicy/%s", &id)
	if count == 0 {
		return []*schema.ResourceData{d}, fmt.Errorf("invalid ID: %s", d.Id())
	}
	_ = d.Set("guid", id)
	return []*schema.ResourceData{d}, nil
}

func ResourceBLRBlobStorePolicy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: importStatePassthroughSetGuidContext,
		},
		CreateContext: resourceBLRBlobStorePolicyCreate,
		ReadContext:   resourceBLRBlobStorePolicyRead,
		DeleteContext: resourceBLRBlobStorePolicyDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"statement": blobStorePolicyStatementSchema(),
			"principal": config.PrincipalSchema(),
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func policyStatementResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"effect": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"action": {
				Type:     schema.TypeSet,
				MaxItems: 4,
				MinItems: 1,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"principal": {
				Type:     schema.TypeSet,
				MinItems: 1,
				MaxItems: 10,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource": {
				Type:     schema.TypeSet,
				MinItems: 1,
				MaxItems: 10,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func blobStorePolicyStatementSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		ForceNew: true,
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return false
		},
		MaxItems: 1,
		Elem:     policyStatementResource(),
	}
}

func schemaToBlobStorePolicy(d *schema.ResourceData) blr.BlobStorePolicy {

	resource := blr.BlobStorePolicy{
		ResourceType: "Bucket",
		Statement:    []blr.BlobStorePolicyStatement{},
	}
	if v, ok := d.GetOk("statement"); ok {
		vL := v.(*schema.Set).List()
		for _, entry := range vL {
			var statement blr.BlobStorePolicyStatement
			mV := entry.(map[string]interface{})
			statement.Effect = mV["effect"].(string)
			statement.Action = tools.ExpandStringList(mV["action"].(*schema.Set).List())
			statement.Principal = tools.ExpandStringList(mV["principal"].(*schema.Set).List())
			statement.Resource = tools.ExpandStringList(mV["resource"].(*schema.Set).List())
			resource.Statement = append(resource.Statement, statement)
		}
	}
	return resource
}

func blobStorePolicyToSchema(resource blr.BlobStorePolicy, d *schema.ResourceData) {
	a := &schema.Set{F: schema.HashResource(policyStatementResource())}
	entry := make(map[string]interface{})
	entry["effect"] = resource.Statement[0].Effect
	entry["action"] = tools.SchemaSetStrings(resource.Statement[0].Action)
	entry["principal"] = tools.SchemaSetStrings(resource.Statement[0].Principal)
	entry["resource"] = tools.SchemaSetStrings(resource.Statement[0].Resource)
	a.Add(entry)

	_ = d.Set("statement", a)
}

func resourceBLRBlobStorePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToBlobStorePolicy(d)

	var created *blr.BlobStorePolicy
	var resp *blr.Response
	err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Configurations.CreateBlobStorePolicy(resource)
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
		return diag.FromErr(fmt.Errorf("failed to create resource: %d", resp.StatusCode()))
	}
	_ = d.Set("guid", created.ID)
	d.SetId(fmt.Sprintf("BlobStorePolicy/%s", created.ID))

	return resourceBLRBucketRead(ctx, d, m)
}

func resourceBLRBlobStorePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "BlobStorePolicy/%s", &id)
	var resource *blr.BlobStorePolicy
	var resp *blr.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Configurations.GetBlobStorePolicyByID(id)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	blobStorePolicyToSchema(*resource, d)
	return diags
}

func resourceBLRBlobStorePolicyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "BlobStorePolicy/%s", &id)
	resource, _, err := client.Configurations.GetBlobStorePolicyByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.Configurations.DeleteBlobStorePolicy(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
