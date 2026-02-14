package iam

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceIAMSMSTemplate() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceIAMSMSTemplateCreate,
		ReadContext:   resourceIAMSMSTemplateRead,
		DeleteContext: resourceIAMSMSTemplateDelete,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
				ForceNew: true,
			},
		},
	}
}

func schemaWriteSMSTemplate(s iam.SMSTemplate, d *schema.ResourceData) error {
	if err := d.Set("organization_id", s.Organization.Value); err != nil {
		return err
	}
	if err := d.Set("message", s.Message); err != nil {
		return err
	}
	if err := d.Set("locale", s.Locale); err != nil {
		return err
	}
	if err := d.Set("external_id", s.ExternalID); err != nil {
		return err
	}
	return nil
}

func schemaReadSMSTemplate(d *schema.ResourceData) (*iam.SMSTemplate, error) {
	var template iam.SMSTemplate

	template.Organization = iam.OrganizationValue{
		Value: d.Get("organization_id").(string),
	}
	template.Locale = d.Get("locale").(string)
	template.ExternalID = d.Get("external_id").(string)
	template.Message = d.Get("message").(string)
	template.Type = d.Get("type").(string)

	return &template, nil
}

func resourceIAMSMSTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := meta.(*config.Config)

	id := d.Id()
	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		_, resp, err = client.SMSTemplates.DeleteSMSTemplate(iam.SMSTemplate{ID: id})
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting SMS template: %w", err))
	}
	d.SetId("")
	return diags
}

func resourceIAMSMSTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var resp *iam.Response
	var template *iam.SMSTemplate

	c := meta.(*config.Config)

	id := d.Id()
	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		template, resp, err = client.SMSTemplates.GetSMSTemplateByID(id)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading SMS template: %w", err))
	}
	// Just in time base64 decoding
	message, err := base64.StdEncoding.DecodeString(template.Message)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error decoding SMS message: %w", err))
	}
	template.Message = string(message)
	if err := schemaWriteSMSTemplate(*template, d); err != nil {
		return diag.FromErr(fmt.Errorf("error setting SMS template: %w", err))
	}
	return diags
}

func resourceIAMSMSTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var createdTemplate *iam.SMSTemplate
	var resp *iam.Response

	c := meta.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	// Refresh token, so we hopefully have SMS permissions to proceed without error
	_ = client.TokenRefresh()

	template, err := schemaReadSMSTemplate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading SMS template: %w", err))
	}
	// Just in time base64 encoding
	template.Message = base64.StdEncoding.EncodeToString([]byte(template.Message))
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		createdTemplate, resp, err = client.SMSTemplates.CreateSMSTemplate(*template)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating SMS template: %w", err))
	}
	d.SetId(createdTemplate.ID)
	return resourceIAMSMSTemplateRead(ctx, d, meta)
}
