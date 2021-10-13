package hsdp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMSMSGatewayConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceIAMSMSGatewayConfigCreate,
		ReadContext:   resourceIAMSMSGatewayConfigRead,
		DeleteContext: resourceIAMSMSGatewayConfigDelete,
		UpdateContext: resourceIAMSMSGatewayConfigUpdate,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_provider": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "twilio",
			},
			"activation_expiry": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  15,
			},
			"properties": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     smsPropertiesSchema(),
			},
			"credentials": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     smsCredentialsSchema(),
			},
			"query_retrieve_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func smsCredentialsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func smsPropertiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from_number": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func schemaWriteSMSGateway(s iam.SMSGateway, d *schema.ResourceData) error {
	if err := d.Set("organization_id", s.Organization.Value); err != nil {
		return err
	}
	if err := d.Set("activation_expiry", s.ActivationExpiry); err != nil {
		return err
	}
	if err := d.Set("gateway_provider", s.Provider); err != nil {
		return err
	}
	properties := make(map[string]interface{})
	properties["sid"] = s.Properties.SID
	properties["endpoint"] = s.Properties.Endpoint
	properties["from_number"] = s.Properties.FromNumber
	p := &schema.Set{F: schema.HashResource(smsPropertiesSchema())}
	p.Add(properties)
	if err := d.Set("properties", p); err != nil {
		return err
	}
	credentials := make(map[string]interface{})
	credentials["token"] = s.Credentials.Token
	c := &schema.Set{F: schema.HashResource(smsCredentialsSchema())}
	c.Add(credentials)
	return d.Set("credentials", c)
}

func schemaReadSMSGateway(d *schema.ResourceData) (*iam.SMSGateway, error) {
	var gw iam.SMSGateway

	gw.Organization = iam.OrganizationValue{
		Value: d.Get("organization_id").(string),
	}
	gw.ActivationExpiry = d.Get("activation_expiry").(int)
	gw.Provider = d.Get("gateway_provider").(string)
	if v, ok := d.GetOk("properties"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			gw.Properties = iam.ProviderProperties{
				SID:        mVi["sid"].(string),
				Endpoint:   mVi["endpoint"].(string),
				FromNumber: mVi["from_number"].(string),
			}
		}
	}
	if v, ok := d.GetOk("credentials"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			gw.Credentials = iam.ProviderCredentials{
				Token: mVi["token"].(string),
			}
		}
	}
	return &gw, nil
}

func resourceIAMSMSGatewayConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var resp *iam.Response
	var gw *iam.SMSGateway
	var err error

	config := meta.(*Config)

	id := d.Id()
	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	serverVersion, _, err := client.SMSGateways.GetSMSGatewayByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	gw, err = schemaReadSMSGateway(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading SMS gateway: %w", err))
	}
	gw.ID = id
	gw.Meta = serverVersion.Meta
	err = tryIAMCall(func() (*iam.Response, error) {
		var err error
		gw, resp, err = client.SMSGateways.UpdateSMSGateway(*gw)
		return resp, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating SMS gateway: %w", err))
	}
	return resourceIAMSMSGatewayConfigRead(ctx, d, meta)
}

func resourceIAMSMSGatewayConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := meta.(*Config)

	id := d.Id()
	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var resp *iam.Response

	err = tryIAMCall(func() (*iam.Response, error) {
		var err error
		_, resp, err = client.SMSGateways.DeleteSMSGateway(iam.SMSGateway{ID: id})
		return resp, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting SMS gateway: %w", err))
	}
	d.SetId("")
	return diags
}

func resourceIAMSMSGatewayConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var resp *iam.Response
	var gw *iam.SMSGateway

	config := meta.(*Config)

	id := d.Id()
	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	err = tryIAMCall(func() (*iam.Response, error) {
		var err error
		gw, resp, err = client.SMSGateways.GetSMSGatewayByID(id)
		return resp, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading SMS gateway: %w", err))
	}
	if err := schemaWriteSMSGateway(*gw, d); err != nil {
		return diag.FromErr(fmt.Errorf("error setting SMS gateway: %w", err))
	}
	return diags
}

func resourceIAMSMSGatewayConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var createdGW *iam.SMSGateway
	var resp *iam.Response

	config := meta.(*Config)

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	// Refresh token, so we hopefully have SMS permissions to proceed without error
	_ = client.TokenRefresh()

	gw, err := schemaReadSMSGateway(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading SMS gateway: %w", err))
	}
	err = tryIAMCall(func() (*iam.Response, error) {
		var err error
		createdGW, resp, err = client.SMSGateways.CreateSMSGateway(*gw)
		return resp, err
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating SMS gateway: %w", err))
	}
	d.SetId(createdGW.ID)
	return resourceIAMSMSGatewayConfigRead(ctx, d, meta)
}
