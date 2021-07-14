package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCDLResearchStudy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLResearchStudyRead,
		Schema: map[string]*schema.Schema{
			"study_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ends_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"study_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceCDLResearchStudyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)
	studyID := d.Get("study_id").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	study, _, err := client.Study.GetStudyByID(studyID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(study.ID)
	_ = d.Set("title", study.Title)
	_ = d.Set("description", study.Description)
	_ = d.Set("ends_at", study.Period.End)
	_ = d.Set("study_owner", study.StudyOwner)

	return diags
}
