package hsdp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/cdl"
	"net/http"
)

func resourceCDLDataTypeDefinition() *schema.Resource {
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
			"jsonschema": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCDLDataTypeDefinitionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("") // This is by design currently
	return diags
}

func resourceCDLDataTypeDefinitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()
	dataTypeDefinition, _, err := client.DataTypeDefinition.GetDataTypeDefinitionByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal([]byte(d.Get("jsonschema").(string)), &((*dataTypeDefinition).JsonSchema))
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
	config := m.(*Config)

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
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
	err = json.Unmarshal([]byte(d.Get("jsonschema").(string)), &dataTypeDefToCreate.JsonSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	createdDtd, resp, err := client.DataTypeDefinition.CreateDataTypeDefinition(dataTypeDefToCreate)
	if err != nil { // currently, creating a DTD with existing name throws 400 and not 409
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusConflict {
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

func resourceCDLDataTypeDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
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

	_ = d.Set("jsonschema", string(b))
	_ = d.Set("createdBy", dataTypeDefinition.CreatedBy)
	_ = d.Set("createdOn", dataTypeDefinition.CreatedOn)
	_ = d.Set("updatedBy", dataTypeDefinition.UpdatedBy)
	_ = d.Set("updatedOn", dataTypeDefinition.UpdatedOn)
	return diags
}
