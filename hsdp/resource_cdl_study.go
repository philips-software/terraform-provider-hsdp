package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/cdl"
)

func resourceCDLStudy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDLStudyCreate,
		ReadContext:   resourceCDLStudyRead,
		UpdateContext: resourceCDLStudyUpdate,
		DeleteContext: resourceCDLStudyDelete,

		Schema: map[string]*schema.Schema{
			"cdl_instance": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"study_owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ends_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCDLStudyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_instance").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	title := d.Get("name").(string)
	description := d.Get("description").(string)
	endsAt := d.Get("ends_at").(string)
	studyOwner := d.Get("study_owner").(string)

	createdStudy, _, err := client.Study.CreateStudy(cdl.Study{
		Title:       title,
		Description: description,
		Period: cdl.Period{
			End: endsAt,
		},
		StudyOwner: studyOwner,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdStudy.ID)
	return diags
}

func resourceCDLStudyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_instance").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()

	study, _, err := client.Study.GetStudyByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", study.Title)
	_ = d.Set("description", study.Description)
	_ = d.Set("study_owner", study.StudyOwner)
	_ = d.Set("ends_at", study.Period.End)
	return diags
}

func resourceCDLStudyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	endpoint := d.Get("cdl_instance").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()
	study, _, err := client.Study.GetStudyByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	study.Title = d.Get("name").(string)
	study.Description = d.Get("description").(string)
	study.Period.End = d.Get("ends_at").(string)
	study.StudyOwner = d.Get("study_owner").(string)

	_, _, err = client.Study.UpdateStudy(*study)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCDLStudyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("") // This is by design currently
	return diags
}
