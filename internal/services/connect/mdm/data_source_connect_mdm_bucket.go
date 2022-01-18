package mdm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceConnectMDMBucket() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectMDMBucketRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_reference_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cors_config_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auditing_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cdn_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"logging_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cross_region_replication_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cache_control_age": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

}

func dataSourceConnectMDMBucketRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	propID := d.Get("proposition_id").(string)
	name := d.Get("name").(string)

	buckets, _, err := client.Buckets.Find(&mdm.GetBucketOptions{
		PropositionID: &propID,
		Name:          &name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if buckets == nil || len(*buckets) == 0 {
		return diag.FromErr(config.ErrResourceNotFound)
	}
	bucket := (*buckets)[0]

	d.SetId(fmt.Sprintf("Bucket/%s", bucket.ID))
	_ = d.Set("guid", bucket.ID)
	if len(bucket.CorsConfiguration) > 0 {
		if corsConfigJSON, err := json.Marshal(bucket.CorsConfiguration); err != nil {
			_ = d.Set("cors_config_json", string(corsConfigJSON))
		}
	}
	if bucket.ReplicationRegionID != nil {
		_ = d.Set("replication_region_id", bucket.ReplicationRegionID.Reference)
	}
	_ = d.Set("description", bucket.Description)
	_ = d.Set("auditing_enabled", bucket.AuditingEnabled)
	_ = d.Set("cdn_enabled", bucket.EnableCDN)
	_ = d.Set("logging_enabled", bucket.LoggingEnabled)
	_ = d.Set("cross_region_replication_enabled", bucket.CrossRegionReplicationEnabled)
	_ = d.Set("cache_control_age", bucket.CacheControlAge)

	return diags
}
