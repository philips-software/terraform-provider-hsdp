package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/connect/mdm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceConnectMDMBucket() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext:      resourceConnectMDMBucketCreate,
		ReadContext:        resourceConnectMDMBucketRead,
		UpdateContext:      resourceConnectMDMBucketUpdate,
		DeleteContext:      resourceConnectMDMBucketDelete,
		DeprecationMessage: "Use the hsdp_blr_bucket resource to manage buckets.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"default_region_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"replication_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cors_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 10,
				Elem:     corsConfigurationsSchema(),
			},
			"versioning_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"logging_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auditing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled_cdn": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cache_control_age": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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

func corsConfigurationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed_origins": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 1,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"allowed_methods": {
				Type:     schema.TypeSet,
				MaxItems: 5,
				MinItems: 1,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"allowed_headers": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_age_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"expose_headers": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func schemaToBucket(d *schema.ResourceData) mdm.Bucket {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	propositionId := d.Get("proposition_id").(string)
	defaultRegionId := d.Get("default_region_id").(string)
	replicationRegionId := d.Get("replication_region_id").(string)
	versioningEnabled := d.Get("versioning_enabled").(bool)
	loggingEnabled := d.Get("logging_enabled").(bool)
	auditingEnabled := d.Get("auditing_enabled").(bool)
	enableCDN := d.Get("enabled_cdn").(bool)
	cacheControlAge := d.Get("cache_control_age").(int)

	resource := mdm.Bucket{
		Name:              name,
		Description:       description,
		PropositionID:     mdm.Reference{Reference: propositionId},
		DefaultRegionID:   mdm.Reference{Reference: defaultRegionId},
		VersioningEnabled: versioningEnabled,
		LoggingEnabled:    loggingEnabled,
		AuditingEnabled:   auditingEnabled,
		EnableCDN:         enableCDN,
		CacheControlAge:   cacheControlAge,
	}
	if replicationRegionId != "" {
		resource.ReplicationRegionID = &mdm.Reference{Reference: replicationRegionId}
	}
	if v, ok := d.GetOk("cors_configuration"); ok {
		vL := v.(*schema.Set).List()
		for _, entry := range vL {
			mV := entry.(map[string]interface{})
			resource.CorsConfiguration = append(resource.CorsConfiguration, mdm.CORSConfiguration{
				MaxAgeSeconds:  mV["max_age_seconds"].(int),
				AllowedOrigins: tools.ExpandStringList(mV["allowed_origins"].(*schema.Set).List()),
				AllowedHeaders: tools.ExpandStringList(mV["allowed_headers"].(*schema.Set).List()),
				AllowedMethods: tools.ExpandStringList(mV["allowed_methods"].(*schema.Set).List()),
				ExposeHeaders:  tools.ExpandStringList(mV["expose_headers"].(*schema.Set).List()),
			})
		}
	}
	return resource
}

func bucketToSchema(resource mdm.Bucket, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("proposition_id", resource.PropositionID)
	_ = d.Set("default_region_id", resource.DefaultRegionID)
	_ = d.Set("guid", resource.ID)
	_ = d.Set("enable_cdn", resource.EnableCDN)
	_ = d.Set("versioning_enabled", resource.VersioningEnabled)
	_ = d.Set("logging_enabled", resource.LoggingEnabled)
	_ = d.Set("auditing_enabled", resource.AuditingEnabled)
	_ = d.Set("cache_control_age", resource.CacheControlAge)
	if resource.ReplicationRegionID != nil {
		_ = d.Set("replication_region_id", resource.ReplicationRegionID.Reference)
	}
	// Add CORSConfiguration
	a := &schema.Set{F: schema.HashResource(corsConfigurationsSchema())}
	for _, cc := range resource.CorsConfiguration {
		entry := make(map[string]interface{})
		entry["allowed_origins"] = tools.SchemaSetStrings(cc.AllowedOrigins)
		entry["allowed_headers"] = tools.SchemaSetStrings(cc.AllowedHeaders)
		entry["expose_headers"] = tools.SchemaSetStrings(cc.ExposeHeaders)
		entry["allowed_methods"] = tools.SchemaSetStrings(cc.AllowedMethods)
		entry["max_age_seconds"] = cc.MaxAgeSeconds
		a.Add(entry)
	}
	_ = d.Set("cors_configuration", a)
}

func resourceConnectMDMBucketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToBucket(d)

	var created *mdm.Bucket
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Buckets.Create(resource)
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
	d.SetId(fmt.Sprintf("Bucket/%s", created.ID))

	return resourceConnectMDMBucketRead(ctx, d, m)
}

func resourceConnectMDMBucketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "Bucket/%s", &id)
	var resource *mdm.Bucket
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Buckets.GetByID(id)
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
	bucketToSchema(*resource, d)
	return diags
}

func resourceConnectMDMBucketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToBucket(d)
	resource.ID = id

	_, _, err = client.Buckets.Update(resource)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceConnectMDMBucketRead(ctx, d, m)
}

func resourceConnectMDMBucketDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.Buckets.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.Buckets.Delete(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
