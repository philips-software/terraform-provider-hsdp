package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceDICOMRepository() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			StateContext: importRepositoryContext,
		},
		CreateContext: resourceDICOMRepositoryCreate,
		ReadContext:   resourceDICOMRepositoryRead,
		DeleteContext: resourceDICOMRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"config_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": { // Query
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository_organization_id": { // Body
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"object_store_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notification": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     notificationSchema(),
			},
			"store_as_composite": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func notificationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDICOMRepositoryDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	operation := func() error {
		var resp *dicom.Response
		_, resp, err = client.Config.DeleteRepository(dicom.Repository{ID: d.Id()}, &dicom.QueryOptions{OrganizationID: &orgID})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceDICOMRepositoryRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	var repo *dicom.Repository
	operation := func() error {
		var resp *dicom.Response
		repo, resp, err = client.Config.GetRepository(d.Id(), &dicom.QueryOptions{OrganizationID: &orgID})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		if errors.Is(err, dicom.ErrNotFound) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	if repo.OrganizationID != orgID {
		_ = d.Set("repository_organization_id", repo.OrganizationID)
	}
	_ = d.Set("object_store_id", repo.ActiveObjectStoreID)
	return diags
}

func resourceDICOMRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	repositoryOrgID := d.Get("repository_organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	repos, _, err := client.Config.GetRepositories(&dicom.QueryOptions{OrganizationID: &orgID})
	if err == nil {
		if len(*repos) > 0 {
			return diag.FromErr(fmt.Errorf("existing dicomRepository found: %s", (*repos)[0].ID))
		}
	}

	defer client.Close()
	repo := dicom.Repository{
		OrganizationID:      orgID,
		ActiveObjectStoreID: d.Get("object_store_id").(string),
	}
	if repositoryOrgID != "" {
		repo.OrganizationID = repositoryOrgID
	}
	if v, ok := d.GetOk("store_as_composite"); ok {
		storeAsComposite := v.(bool)
		repo.StoreAsComposite = &storeAsComposite
	}
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
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.ID)
	return resourceDICOMRepositoryRead(ctx, d, m)
}
