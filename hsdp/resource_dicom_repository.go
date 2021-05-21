package hsdp

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
)

func resourceDICOMRepository() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
		},
	}
}

func resourceDICOMRepositoryDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	operation := func() error {
		var resp *dicom.Response
		_, resp, err = client.Config.DeleteRepository(dicom.Repository{ID: d.Id()}, &dicom.QueryOptions{OrganizationID: &orgID})
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceDICOMRepositoryRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	var repo *dicom.Repository
	operation := func() error {
		var resp *dicom.Response
		repo, resp, err = client.Config.GetRepository(d.Id(), &dicom.QueryOptions{OrganizationID: &orgID})
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return diag.FromErr(err)
	}
	if repo.OrganizationID != orgID {
		_ = d.Set("repository_organization_id", repo.OrganizationID)
	}
	_ = d.Set("object_store_id", repo.ActiveObjectStoreID)
	return diags
}

func resourceDICOMRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	repositoryOrgID := d.Get("repository_organization_id").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	repo := dicom.Repository{
		OrganizationID:      orgID,
		ActiveObjectStoreID: d.Get("object_store_id").(string),
	}
	if repositoryOrgID != "" {
		repo.OrganizationID = repositoryOrgID
	}
	var created *dicom.Repository
	operation := func() error {
		var resp *dicom.Response
		created, resp, err = client.Config.CreateRepository(repo, &dicom.QueryOptions{OrganizationID: &orgID})
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.ID)
	return resourceDICOMRepositoryRead(ctx, d, m)
}
