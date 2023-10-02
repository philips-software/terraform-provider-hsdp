package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMBlobDataContract() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceConnectMDMBlobDataContractCreate,
		ReadContext:   resourceConnectMDMBlobDataContractRead,
		UpdateContext: resourceConnectMDMBlobDataContractUpdate,
		DeleteContext: resourceConnectMDMBlobDataContractDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"data_type_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"bucket_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"storage_class_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"root_path_in_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logging_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"cross_region_replication_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaToBlobDataContract(d *schema.ResourceData) mdm.BlobDataContract {
	name := d.Get("name").(string)
	dataTypeId := d.Get("data_type_id").(string)
	bucketId := d.Get("bucket_id").(string)
	storageClassId := d.Get("storage_class_id").(string)
	rootPathInBucket := d.Get("root_path_in_bucket").(string)
	loggingEnabled := d.Get("logging_enabled").(bool)
	crossRegionReplicationEnabled := d.Get("cross_region_replication_enabled").(bool)

	resource := mdm.BlobDataContract{
		Name:                          name,
		DataTypeID:                    mdm.Reference{Reference: dataTypeId},
		BucketID:                      mdm.Reference{Reference: bucketId},
		StorageClassID:                mdm.Reference{Reference: storageClassId},
		RootPathInBucket:              rootPathInBucket,
		LoggingEnabled:                loggingEnabled,
		CrossRegionReplicationEnabled: crossRegionReplicationEnabled,
	}
	return resource
}

func blobDataContractToSchema(resource mdm.BlobDataContract, d *schema.ResourceData) {
	_ = d.Set("root_path_in_bucket", resource.RootPathInBucket)
	_ = d.Set("name", resource.Name)
	_ = d.Set("data_type_id", resource.DataTypeID.Reference)
	_ = d.Set("bucket_id", resource.BucketID.Reference)
	_ = d.Set("storage_class_id", resource.StorageClassID.Reference)
	_ = d.Set("guid", resource.ID)
	_ = d.Set("logging_enabled", resource.LoggingEnabled)
	_ = d.Set("cross_region_replication_enabled", resource.CrossRegionReplicationEnabled)
}

func resourceConnectMDMBlobDataContractCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToBlobDataContract(d)

	var created *mdm.BlobDataContract
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.BlobDataContracts.Create(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if created == nil {
		return diag.FromErr(fmt.Errorf("failed to create resource: %d", resp.StatusCode()))
	}
	_ = d.Set("guid", created.ID)
	d.SetId(fmt.Sprintf("BlobDataContract/%s", created.ID))
	return resourceConnectMDMBlobDataContractRead(ctx, d, m)
}

func resourceConnectMDMBlobDataContractRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "BlobDataContract/%s", &id)
	var resource *mdm.BlobDataContract
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.BlobDataContracts.GetByID(id)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	blobDataContractToSchema(*resource, d)
	return diags
}

func resourceConnectMDMBlobDataContractUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToBlobDataContract(d)
	resource.ID = id

	_, _, err = client.BlobDataContracts.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMBlobDataContractRead(ctx, d, m)
}

func resourceConnectMDMBlobDataContractDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.BlobDataContracts.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.BlobDataContracts.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
