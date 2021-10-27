package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceDockerRegistry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDockerRegistryRead,
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"total_pulls": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_tags": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"updated_at": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_digests": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_pulls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"compressed_sizes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}

}

func dataSourceDockerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespaceId := d.Get("namespace_id").(string)
	name := d.Get("name").(string)

	repo, err := client.Repositories.GetRepository(ctx, namespaceId, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading repository: %w", err))
	}
	tagList, err := client.Repositories.GetTags(ctx, repo.ID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading tags: %w", err))
	}

	var ids []int
	var tags []string
	var updatedAt []string
	var imageIDs []string
	var imageDigests []string
	var compressedSizes []int
	var numPulls []int

	for _, t := range *tagList {
		ids = append(ids, t.ID)
		tags = append(tags, t.Name)
		updatedAt = append(updatedAt, t.UpdatedAt.Format(time.RFC3339))
		imageIDs = append(imageIDs, t.ImageId)
		imageDigests = append(imageDigests, t.Digest)
		compressedSizes = append(compressedSizes, t.CompressedSize)
		numPulls = append(numPulls, t.NumPulls)
	}
	_ = d.Set("ids", ids)
	_ = d.Set("tags", tags)
	_ = d.Set("updated_at", updatedAt)
	_ = d.Set("image_ids", imageIDs)
	_ = d.Set("image_digests", imageDigests)
	_ = d.Set("num_pulls", numPulls)
	_ = d.Set("compressed_sizes", compressedSizes)
	_ = d.Set("total_pulls", repo.NumPulls)
	_ = d.Set("total_tags", repo.NumTags)
	d.SetId(repo.ID)
	return diags
}
