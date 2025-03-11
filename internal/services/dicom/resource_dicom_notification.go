package dicom

import (
	"context"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/dip-software/go-dip-api/dicom"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceDICOMNotification() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDICOMNotificationCreate,
		ReadContext:   resourceDICOMNotificationRead,
		DeleteContext: resourceDICOMNotificationDelete,

		Schema: map[string]*schema.Schema{
			"config_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "this field should be removed as it's no longer used by this resource",
				ForceNew:   true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"endpoint_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_organization_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDICOMNotificationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	var notification *dicom.Notification
	var resp *dicom.Response
	operation := func() error {
		notification, resp, err = client.Config.GetNotification(&dicom.QueryOptions{})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diags
		}
	}
	notification.Enabled = false
	notification.ID = ""
	_, _, _ = client.Config.CreateNotification(*notification, &dicom.QueryOptions{})
	d.SetId("")
	return diags
}

func resourceDICOMNotificationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	var notification *dicom.Notification
	operation := func() error {
		var resp *dicom.Response
		notification, resp, err = client.Config.GetNotification(&dicom.QueryOptions{})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil { // For now just declare the notification not there in case of error
		d.SetId("")
		return diags
	}
	_ = d.Set("endpoint_url", notification.Endpoint)
	_ = d.Set("default_organization_id", notification.DefaultOrganizationID)
	_ = d.Set("enabled", notification.Enabled)
	return diags
}

func resourceDICOMNotificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	endpointURL := d.Get("endpoint_url").(string)
	defaultOrganizationID := d.Get("default_organization_id").(string)
	enabled := d.Get("enabled").(bool)

	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	resource := dicom.Notification{
		Enabled:               enabled,
		Endpoint:              endpointURL,
		DefaultOrganizationID: defaultOrganizationID,
	}

	var created *dicom.Notification
	operation := func() error {
		var resp *dicom.Response
		created, resp, err = client.Config.CreateNotification(resource, &dicom.QueryOptions{})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.ID)
	return resourceDICOMNotificationRead(ctx, d, m)
}
