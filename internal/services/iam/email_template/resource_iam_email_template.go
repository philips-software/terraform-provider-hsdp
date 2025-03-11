package email_template

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

var descriptions = map[string]string{
	"email_template": "There are certain IAM flows that trigger email notifications to the user. The email delivered to users will use the IAM specific email template by default. This resource allows you to manage template for your HSDP IAM organization.",
}

func ResourceIAMEmailTemplate() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["email_template"],
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceIAMEmailTemplateCreate,
		ReadContext:   resourceIAMEmailTemplateRead,
		DeleteContext: resourceIAMEmailTemplateDelete,

		Schema: map[string]*schema.Schema{
			"managing_organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Id of the IAM Org to apply this email template to.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The email template type.",
			},
			"from": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressDefault,
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "HTML",
				Description: "The template format. Must be 'HTML' currently.",
			},
			"subject": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
				Description: "The Subject line of the email.",
			},
			"message": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The message body.",
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressDefault,
					tools.SuppressCaseDiffs,
				),
				Description: "The locale of the template. When not specified the template will become the default. Only a single default template is allowed of course.",
			},
			"link": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressDefault,
				Description:      "A clickable link, depends on the template type.",
			},
			"message_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMEmailTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var template iam.EmailTemplate

	template.Type = d.Get("type").(string)
	template.Format = d.Get("format").(string)
	template.Subject = d.Get("subject").(string)
	template.Message = base64.StdEncoding.EncodeToString([]byte(d.Get("message").(string)))
	template.Link = d.Get("link").(string)
	template.Locale = d.Get("locale").(string)
	template.From = d.Get("from").(string)
	template.ManagingOrganization = d.Get("managing_organization").(string)

	var createdTemplate *iam.EmailTemplate
	var resp *iam.Response
	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		createdTemplate, resp, err = client.EmailTemplates.CreateTemplate(template)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})

	if err != nil {
		if resp.StatusCode() == http.StatusConflict {
			templates, _, getErr := client.EmailTemplates.GetTemplates(&iam.GetEmailTemplatesOptions{
				Type:           &template.Type,
				OrganizationID: &template.ManagingOrganization,
				Locale:         &template.Locale,
			})
			if getErr != nil {
				return diag.FromErr(fmt.Errorf("createEmailTemplate HTTP 409 conflict: %w", getErr))
			}
			if len(*templates) > 0 {
				return diag.FromErr(fmt.Errorf("conflicting template with ID '%s': %w", (*templates)[0].ID, err))
			}
		}
		return diag.FromErr(err)
	}
	_ = d.Set("message_base64", createdTemplate.Message)
	d.SetId(createdTemplate.ID)
	return resourceIAMEmailTemplateRead(ctx, d, m)
}

func resourceIAMEmailTemplateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	template, _, err := client.EmailTemplates.GetTemplateByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("subject", template.Subject)
	if template.Locale != "default" {
		_ = d.Set("locale", template.Locale)
	}
	_ = d.Set("from", template.From)
	_ = d.Set("format", template.Format)
	_ = d.Set("type", template.Type)
	_ = d.Set("link", template.Link)
	// Message is not returned in the read call

	return diags
}

func resourceIAMEmailTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var template iam.EmailTemplate
	template.ID = d.Id()
	var ok bool
	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *iam.Response
		var err error
		ok, resp, err = client.EmailTemplates.DeleteTemplate(template)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrDeleteEmailTemplateFailed)
	}
	d.SetId("")
	return diags
}
