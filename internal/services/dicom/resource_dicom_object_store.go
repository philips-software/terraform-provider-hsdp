package dicom

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/dip-software/go-dip-api/dicom"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceDICOMObjectStore() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDICOMObjectStoreCreate,
		ReadContext:   resourceDICOMObjectStoreRead,
		DeleteContext: resourceDICOMObjectStoreDelete,
		SchemaVersion: 1,

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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"static_access": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     staticAccessSchema(),
			},
			"access_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3creds_access": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     s3credsAccessSchema(),
			},
		},
	}
}

func s3credsAccessSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"folder_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"service_account": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     serviceAccountSchema(),
			},
		},
	}
}

func staticAccessSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
		},
	}
}

func serviceAccountSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"access_token_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDICOMObjectStoreDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	forceDelete := d.Get("force_delete").(bool)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	if !forceDelete { // soft delete
		d.SetId("")
		return diags
	}
	operation := func() error {
		var resp *dicom.Response
		_, resp, err = client.Config.DeleteObjectStore(dicom.ObjectStore{ID: d.Id()}, &dicom.QueryOptions{
			OrganizationID: &orgID,
		})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceDICOMObjectStoreRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	orgID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	var store *dicom.ObjectStore
	operation := func() error {
		var resp *dicom.Response
		store, resp, err = client.Config.GetObjectStore(d.Id(), &dicom.QueryOptions{
			OrganizationID: &orgID,
		})
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
	_ = d.Set("description", store.Description)
	_ = d.Set("access_type", store.AccessType)
	if store.StaticAccess != nil {
		staticSettings := make(map[string]interface{})
		staticSettings["endpoint"] = store.StaticAccess.Endpoint
		staticSettings["bucket_name"] = store.StaticAccess.BucketName
		staticSettings["access_key"] = store.StaticAccess.AccessKey
		staticSettings["secret_key"] = store.StaticAccess.SecretKey
		s := &schema.Set{F: schema.HashResource(staticAccessSchema())}
		s.Add(staticSettings)
		_ = d.Set("static_access", s)
	}
	if store.CredServiceAccess != nil {
		credsSettings := make(map[string]interface{})
		credsSettings["endpoint"] = store.CredServiceAccess.Endpoint
		credsSettings["folder_path"] = store.CredServiceAccess.FolderPath
		credsSettings["bucket_name"] = store.CredServiceAccess.BucketName
		credsSettings["product_key"] = store.CredServiceAccess.ProductKey

		accountSettings := make(map[string]interface{})
		accountSettings["service_id"] = store.CredServiceAccess.ServiceAccount.ServiceID
		accountSettings["private_key"] = getS3CredsPrivateKeyFromState(d)
		accountSettings["access_token_endpoint"] = store.CredServiceAccess.ServiceAccount.AccessTokenEndPoint
		accountSettings["token_endpoint"] = store.CredServiceAccess.ServiceAccount.TokenEndPoint
		accountSettings["name"] = store.CredServiceAccess.ServiceAccount.Name
		s := &schema.Set{F: schema.HashResource(serviceAccountSchema())}
		s.Add(accountSettings)
		credsSettings["service_account"] = s

		c := &schema.Set{F: schema.HashResource(s3credsAccessSchema())}
		c.Add(credsSettings)
		_ = d.Set("s3creds_access", c)
	}

	return diags
}

func resourceDICOMObjectStoreCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	orgID := d.Get("organization_id").(string)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	store := dicom.ObjectStore{}
	store.Description = d.Get("description").(string)

	if v, ok := d.GetOk("static_access"); ok {
		vL := v.(*schema.Set).List()
		staticAccess := &dicom.StaticAccess{}
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			staticAccess.Endpoint = mVi["endpoint"].(string)
			staticAccess.BucketName = mVi["bucket_name"].(string)
			staticAccess.AccessKey = mVi["access_key"].(string)
			staticAccess.SecretKey = mVi["secret_key"].(string)
		}
		store.StaticAccess = staticAccess
		store.AccessType = "static"
	}
	if v, ok := d.GetOk("s3creds_access"); ok {
		vL := v.(*schema.Set).List()
		credsAccess := &dicom.CredsServiceAccess{}
		var aVi []interface{}
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			credsAccess.Endpoint = mVi["endpoint"].(string)
			credsAccess.BucketName = mVi["bucket_name"].(string)
			credsAccess.FolderPath = mVi["folder_path"].(string)
			credsAccess.ProductKey = mVi["product_key"].(string)
			aVi = mVi["service_account"].(*schema.Set).List()
		}
		for _, vi := range aVi {
			mVi := vi.(map[string]interface{})
			credsAccess.ServiceAccount.AccessTokenEndPoint = mVi["access_token_endpoint"].(string)
			credsAccess.ServiceAccount.PrivateKey = mVi["private_key"].(string)
			credsAccess.ServiceAccount.ServiceID = mVi["service_id"].(string)
			credsAccess.ServiceAccount.TokenEndPoint = mVi["token_endpoint"].(string)
			credsAccess.ServiceAccount.Name = mVi["name"].(string)
		}
		store.CredServiceAccess = credsAccess
		store.AccessType = "s3Creds"
	}
	var created *dicom.ObjectStore
	operation := func() error {
		var resp *dicom.Response
		created, resp, err = client.Config.CreateObjectStore(store, &dicom.QueryOptions{
			OrganizationID: &orgID,
		})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	if created == nil || created.ID == "" {
		return diag.FromErr(fmt.Errorf("failed to create object store, even though no error was reported"))
	}
	d.SetId(created.ID)
	return resourceDICOMObjectStoreRead(ctx, d, m)
}

func getS3CredsPrivateKeyFromState(d *schema.ResourceData) string {
	privateKey := ""
	if v, ok := d.GetOk("s3creds_access"); ok {
		vL := v.(*schema.Set).List()
		var aVi []interface{}
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			aVi = mVi["service_account"].(*schema.Set).List()
		}
		for _, vi := range aVi {
			mVi := vi.(map[string]interface{})
			privateKey = mVi["private_key"].(string)
		}
	}
	return privateKey
}
