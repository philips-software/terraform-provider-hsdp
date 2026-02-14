package email_template

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/philips-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMEmailTemplates() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["email_template"],
		ReadContext: dataSourceIAMEmailTemplates,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"formats": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"locales": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subjects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"messages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"from": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceIAMEmailTemplates(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := d.Get("organization_id").(string)
	var locale *string
	var templateType *string

	if val, ok := d.GetOk("locale"); ok {
		s := val.(string)
		locale = &s
	}
	if val, ok := d.GetOk("type"); ok {
		s := val.(string)
		templateType = &s
	}
	resources, _, err := client.EmailTemplates.GetTemplates(&iam.GetEmailTemplatesOptions{
		OrganizationID: &orgID,
		Locale:         locale,
		Type:           templateType,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	var subjects []string
	var messages []string
	var from []string
	var links []string
	var types []string
	var formats []string
	var locales []string

	for _, resource := range *resources {
		ids = append(ids, resource.ID)
		subjects = append(subjects, resource.Subject)
		decodedMessage, err := base64.StdEncoding.DecodeString(resource.Message)
		if err != nil {
			messages = append(messages, fmt.Sprintf("error decoding: %v", err))
		} else {
			messages = append(messages, string(decodedMessage))
		}
		from = append(from, resource.From)
		links = append(links, resource.Link)
		types = append(types, resource.Type)
		formats = append(formats, resource.Format)
		locales = append(locales, resource.Locale)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("subjects", subjects)
	_ = d.Set("messages", messages)
	_ = d.Set("from", from)
	_ = d.Set("links", links)
	_ = d.Set("types", types)
	_ = d.Set("formats", formats)
	_ = d.Set("locales", locales)

	d.SetId(orgID)
	return diags
}
