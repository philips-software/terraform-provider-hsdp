package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/philips-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMDataType() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMDataTypeCreate,
		ReadContext:   resourceConnectMDMDataTypeRead,
		UpdateContext: resourceConnectMDMDataTypeUpdate,
		DeleteContext: resourceConnectMDMDataTypeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaToDataType(d *schema.ResourceData) mdm.DataType {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	propositionId := d.Get("proposition_id").(string)
	tags := tools.ExpandStringList(d.Get("tags").(*schema.Set).List())

	resource := mdm.DataType{
		Name:          name,
		Description:   description,
		PropositionId: mdm.Reference{Reference: propositionId},
		Tags:          tags,
	}
	return resource
}

func dataTypeToSchema(resource mdm.DataType, d *schema.ResourceData) {
	_ = d.Set("name", resource.Name)
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("tags", tools.SchemaSetStrings(resource.Tags))
	_ = d.Set("guid", resource.ID)
	_ = d.Set("proposition_id", resource.PropositionId.Reference)
}

func resourceConnectMDMDataTypeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToDataType(d)

	var created *mdm.DataType
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.DataTypes.Create(resource)
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
	d.SetId(fmt.Sprintf("DataType/%s", created.ID))
	return resourceConnectMDMDataTypeRead(ctx, d, m)
}

func resourceConnectMDMDataTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "DataType/%s", &id)
	var resource *mdm.DataType
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.DataTypes.GetByID(id)
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
	dataTypeToSchema(*resource, d)
	return diags
}

func resourceConnectMDMDataTypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	service := schemaToDataType(d)
	service.ID = id

	_, _, err = client.DataTypes.Update(service)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMDataTypeRead(ctx, d, m)
}

func resourceConnectMDMDataTypeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.DataTypes.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.DataTypes.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
