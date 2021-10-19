package cdl

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceCDLLabelDefinition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLLabelDefinitionRead,
		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label_def_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"study_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label_def_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceCDLLabelDefinitionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	labelDefId := d.Get("label_def_id").(string)
	endpoint := d.Get("cdl_endpoint").(string)
	studyId := d.Get("study_id").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	labelDefinition, _, err := client.LabelDefinition.GetLabelDefinitionByID(studyId, labelDefId)
	if err != nil {
		return diag.FromErr(err)
	}

	labelScopeBytes, err := json.Marshal((*labelDefinition).LabelScope)
	if err != nil {
		return diag.FromErr(err)
	}

	labelsBytes, err := json.Marshal((*labelDefinition).Labels)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId((*labelDefinition).ID)
	_ = d.Set("label_def_name", (*labelDefinition).LabelDefName)
	_ = d.Set("description", (*labelDefinition).Description)
	_ = d.Set("label_scope", string(labelScopeBytes))
	_ = d.Set("label_name", (*labelDefinition).Label)
	_ = d.Set("type", (*labelDefinition).Type)
	_ = d.Set("labels", string(labelsBytes))
	_ = d.Set("created_by", (*labelDefinition).CreatedBy)
	_ = d.Set("created_on", (*labelDefinition).CreatedOn)

	return diags
}
