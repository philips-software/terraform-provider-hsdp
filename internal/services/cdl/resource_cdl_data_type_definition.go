package cdl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/cdl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceCDLDataTypeDefinition() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDLDataTypeDefinitionCreate,
		ReadContext:   resourceCDLDataTypeDefinitionRead,
		UpdateContext: resourceCDLDataTypeDefinitionUpdate,
		DeleteContext: resourceCDLDataTypeDefinitionDelete,

		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"json_schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCDLDataTypeDefinitionDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("") // This is by design currently
	return diags
}

func resourceCDLDataTypeDefinitionUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()
	dataTypeDefinition, _, err := client.DataTypeDefinition.GetDataTypeDefinitionByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal([]byte(d.Get("json_schema").(string)), &((*dataTypeDefinition).JsonSchema))
	if err != nil {
		return diag.FromErr(err)
	}
	(*dataTypeDefinition).Description = d.Get("description").(string)

	_, _, err = client.DataTypeDefinition.UpdateDataTypeDefinition(*dataTypeDefinition)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCDLDataTypeDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	dtdName := d.Get("name").(string)
	dtdDescription := d.Get("description").(string)
	dataTypeDefToCreate := cdl.DataTypeDefinition{
		Name:        dtdName,
		Description: dtdDescription,
	}
	err = json.Unmarshal([]byte(d.Get("json_schema").(string)), &dataTypeDefToCreate.JsonSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	createdDtd, resp, err := client.DataTypeDefinition.CreateDataTypeDefinition(dataTypeDefToCreate)
	if err != nil { // currently, creating a DTD with existing name throws 400 and not 409
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode() != http.StatusConflict {
			return diag.FromErr(err)
		}
		// Search for existing DTD
		dataTypeDefinitions, _, err2 := client.DataTypeDefinition.GetDataTypeDefinitions(nil)
		if err2 != nil {
			return diag.FromErr(fmt.Errorf("on match attempt during Create conflict: %w", err))
		}
		for _, dtd := range dataTypeDefinitions {
			if dtd.Name == dtdName {
				d.SetId(dtd.ID)
				return resourceCDLDataTypeDefinitionRead(ctx, d, m)
			}
		}
		return diag.FromErr(err)
	}
	d.SetId(createdDtd.ID)
	return resourceCDLDataTypeDefinitionRead(ctx, d, m)
}

func resourceCDLDataTypeDefinitionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()

	dataTypeDefinition, _, err := client.DataTypeDefinition.GetDataTypeDefinitionByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", dataTypeDefinition.Name)
	_ = d.Set("description", dataTypeDefinition.Description)

	b, err := json.Marshal(dataTypeDefinition.JsonSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("json_schema", string(b))
	_ = d.Set("created_by", dataTypeDefinition.CreatedBy)
	_ = d.Set("created_on", dataTypeDefinition.CreatedOn)
	_ = d.Set("updated_by", dataTypeDefinition.UpdatedBy)
	_ = d.Set("updated_on", dataTypeDefinition.UpdatedOn)
	return diags
}
