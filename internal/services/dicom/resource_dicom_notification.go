package dicom

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceDICOMNotification() *schema.Resource {
	return &schema.Resource{
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
				Type:     schema.TypeString,
				Required: true,
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

func resourceDICOMNotificationDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// This is a NOOP for now
	d.SetId("")
	return diags
}

func resourceDICOMNotificationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	var notification *dicom.Notification
	operation := func() error {
		var resp *dicom.Response
		notification, resp, err = client.Config.GetNotification(&dicom.QueryOptions{OrganizationID: &orgID})
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("endpoint_url", notification.Endpoint)
	_ = d.Set("default_organization_id", notification.DefaultOrganizationID)
	return diags
}

func resourceDICOMNotificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	repo := dicom.Repository{
		OrganizationID:      orgID,
		ActiveObjectStoreID: d.Get("object_store_id").(string),
	}
	repo.OrganizationID = orgID
	if v, ok := d.GetOk("notification"); ok {
		vL := v.(*schema.Set).List()
		repoNotification := dicom.RepositoryNotification{}
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			repoNotification.Enabled = mVi["enabled"].(bool)
			repoNotification.OrganizationID = mVi["organization_id"].(string)
		}
		repo.Notification = &repoNotification
	}

	var created *dicom.Repository
	operation := func() error {
		var resp *dicom.Response
		created, resp, err = client.Config.CreateRepository(repo, &dicom.QueryOptions{OrganizationID: &orgID})
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.ID)
	return resourceDICOMNotificationRead(ctx, d, m)
}
