package hsdp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/cdl"
)

func resourceCDLLabelDefinition() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDLLabelDefinitionCreate,
		ReadContext:   resourceCDLLabelDefinitionRead,
		UpdateContext: resourceCDLLabelDefinitionUpdate,
		DeleteContext: resourceCDLLabelDefinitionDelete,

		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"study_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label_def_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label_scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceCDLLabelDefinitionUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceCDLLabelDefinitionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)

	endpoint := d.Get("cdl_endpoint").(string)
	study_id := d.Get("study_id").(string)
	label_def_id := d.Id()

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	resp, err := client.LabelDefinition.DeleteLabelDefinitionById(study_id, label_def_id)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}
	return diags
}

func resourceCDLLabelDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	endpoint := d.Get("cdl_endpoint").(string)
	study_id := d.Get("study_id").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	labelDefName := d.Get("label_def_name").(string)
	labelDefDescription := d.Get("description").(string)
	label := d.Get("label_name").(string)
	labelType := d.Get("type").(string)

	labelDefToCreate := cdl.LabelDefinition{
		LabelDefName: labelDefName,
		Description:  labelDefDescription,
		Label:        label,
		Type:         labelType,
		LabelScope: cdl.LabelScope{
			Type: d.Get("label_scope").(string),
		},
	}

	labelsArray := expandStringList(d.Get("labels").(*schema.Set).List())
	for _, l := range labelsArray {
		labelArrayElem := cdl.LabelsArrayElem{
			Label: l,
		}
		labelDefToCreate.Labels = append(labelDefToCreate.Labels, labelArrayElem)
	}

	createdLabelDef, resp, err := client.LabelDefinition.CreateLabelDefinition(study_id, labelDefToCreate)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusConflict {
			return diag.FromErr(err)
		}
		// Search for existing Label def
		createdLabelDefs, _, err2 := client.LabelDefinition.GetLabelDefinitions(study_id, &cdl.GetOptions{})
		if err2 != nil {
			return diag.FromErr(fmt.Errorf("on match attempt during Create conflict: %w", err))
		}
		for _, labelDef := range createdLabelDefs {
			if labelDef.LabelDefName == labelDefName {
				d.SetId(labelDef.ID)
				return resourceCDLLabelDefinitionRead(ctx, d, m)
			}
		}
		return diag.FromErr(err)
	}
	d.SetId(createdLabelDef.ID)
	return resourceCDLLabelDefinitionRead(ctx, d, m)
}

func resourceCDLLabelDefinitionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)
	study_id := d.Get("study_id").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()

	labelDefinition, _, err := client.LabelDefinition.GetLabelDefinitionByID(study_id, id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("label_def_name", labelDefinition.LabelDefName)
	_ = d.Set("description", labelDefinition.Description)

	labelScopeBytes, err := json.Marshal(labelDefinition.LabelScope)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("label_scope", string(labelScopeBytes))

	_ = d.Set("label_name", labelDefinition.Label)
	_ = d.Set("type", labelDefinition.Type)

	var labelsArray []string
	for _, l := range labelDefinition.Labels {
		labelsArray = append(labelsArray, l.Label)
	}
	_ = d.Set("labels", labelsArray)
	_ = d.Set("created_by", labelDefinition.CreatedBy)
	_ = d.Set("created_on", labelDefinition.CreatedOn)
	return diags
}
