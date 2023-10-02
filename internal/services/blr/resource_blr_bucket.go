package blr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/philips-software/go-hsdp-api/blr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceBLRBucket() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceBLRBucketCreate,
		ReadContext:   resourceBLRBucketRead,
		UpdateContext: resourceBLRBucketUpdate,
		DeleteContext: resourceBLRBucketDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"principal": config.PrincipalSchema(),
			"price_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cors_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     corsConfigurationsSchema(),
			},
			"enable_create_or_delete_blob_meta": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_hsdp_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_cdn": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cache_control_age": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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
				Default:  1,
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

func schemaToBucket(d *schema.ResourceData) blr.Bucket {
	name := d.Get("name").(string)
	propositionId := d.Get("proposition_id").(string)
	priceClass := d.Get("price_class").(string)
	enableHSDPDomain := d.Get("enable_hsdp_domain").(bool)
	enableCreateOrDeleteBlobMeta := d.Get("enable_create_or_delete_blob_meta").(bool)
	enableCDN := d.Get("enable_cdn").(bool)
	cacheControlAge := d.Get("cache_control_age").(int)

	resource := blr.Bucket{
		ResourceType:                 "Bucket",
		Name:                         name,
		PropositionID:                blr.Reference{Reference: propositionId, Display: "Terraform managed"},
		PriceClass:                   priceClass,
		EnableHSDPDomain:             enableHSDPDomain,
		EnableCreateOrDeleteBlobMeta: enableCreateOrDeleteBlobMeta,
		EnableCDN:                    enableCDN,
		CacheControlAge:              cacheControlAge,
	}
	if v, ok := d.GetOk("cors_configuration"); ok {
		vL := v.(*schema.Set).List()
		for _, entry := range vL {
			mV := entry.(map[string]interface{})
			resource.CorsConfiguration.MaxAgeSeconds = mV["max_age_seconds"].(int)
			resource.CorsConfiguration.AllowedOrigins = tools.ExpandStringList(mV["allowed_origins"].(*schema.Set).List())
			resource.CorsConfiguration.AllowedHeaders = tools.ExpandStringList(mV["allowed_headers"].(*schema.Set).List())
			resource.CorsConfiguration.AllowedMethods = tools.ExpandStringList(mV["allowed_methods"].(*schema.Set).List())
			resource.CorsConfiguration.ExposeHeaders = tools.ExpandStringList(mV["expose_headers"].(*schema.Set).List())
		}
	}
	return resource
}

func bucketToSchema(resource blr.Bucket, d *schema.ResourceData) {
	_ = d.Set("name", resource.Name)
	_ = d.Set("proposition_id", resource.PropositionID.Reference)
	_ = d.Set("enable_cdn", resource.EnableCDN)
	_ = d.Set("price_class", resource.PriceClass)
	_ = d.Set("enable_hsdp_domain", resource.EnableHSDPDomain)
	_ = d.Set("enable_create_or_delete_blob_meta", resource.EnableCreateOrDeleteBlobMeta)
	_ = d.Set("cache_control_age", resource.CacheControlAge)

	// Add CORSConfiguration
	a := &schema.Set{F: schema.HashResource(corsConfigurationsSchema())}
	entry := make(map[string]interface{})
	entry["allowed_origins"] = tools.SchemaSetStrings(resource.CorsConfiguration.AllowedOrigins)
	entry["allowed_headers"] = tools.SchemaSetStrings(resource.CorsConfiguration.AllowedHeaders)
	entry["expose_headers"] = tools.SchemaSetStrings(resource.CorsConfiguration.ExposeHeaders)
	entry["allowed_methods"] = tools.SchemaSetStrings(resource.CorsConfiguration.AllowedMethods)
	entry["max_age_seconds"] = resource.CorsConfiguration.MaxAgeSeconds
	a.Add(entry)

	_ = d.Set("cors_configuration", a)
}

func resourceBLRBucketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToBucket(d)

	var created *blr.Bucket
	var resp *blr.Response
	err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Configurations.CreateBucket(resource)
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

	return resourceBLRBucketRead(ctx, d, m)
}

func resourceBLRBucketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "Bucket/%s", &id)
	var resource *blr.Bucket
	var resp *blr.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Configurations.GetBucketByID(id)
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

func resourceBLRBucketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)

	resource := schemaToBucket(d)
	resource.ID = id

	_, _, err = client.Configurations.UpdateBucket(resource)

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		return diags
	}
	return resourceBLRBucketRead(ctx, d, m)
}

func resourceBLRBucketDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	principal := config.SchemaToPrincipal(d, m)

	client, err := c.BLRClient(principal)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("guid").(string)
	resource, _, err := client.Configurations.GetBucketByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.Configurations.DeleteBucket(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	d.SetId("")
	return diags
}
