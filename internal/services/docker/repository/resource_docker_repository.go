package repository

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/console/docker"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceDockerRepository() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDockerRepositoryCreate,
		ReadContext:   resourceDockerRepositoryRead,
		UpdateContext: resourceDockerRepositoryUpdate,
		DeleteContext: resourceDockerRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"short_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_pulls": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_tags": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"latest_tag": {
				Type:     schema.TypeString,
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

func resourceDockerRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	shortDescription := d.Get("short_description").(string)
	fullDescription := d.Get("full_description").(string)

	_, err = client.Repositories.UpdateRepository(ctx, docker.Repository{ID: id}, docker.RepositoryDetailsInput{
		ShortDescription: shortDescription,
		FullDescription:  fullDescription,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("updateRepository: %w", err))
	}

	return dataSourceDockerRepositoryRead(ctx, d, m)
}

func resourceDockerRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	err = client.Repositories.DeleteRepository(ctx, docker.Repository{ID: id})
	if err != nil {
		return diag.FromErr(fmt.Errorf("deleteRepository: %w", err))
	}
	d.SetId("")
	return diags
}

func resourceDockerRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	readDiags := dataSourceDockerRepositoryRead(ctx, d, m)
	if len(readDiags) > 0 { // Currently, the only way to discover removed repositories
		d.SetId("")
	}
	return diags
}

func resourceDockerRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)
	namespaceId := d.Get("namespace_id").(string)
	shortDescription := d.Get("short_description").(string)
	fullDescription := d.Get("full_description").(string)

	repo, err := client.Repositories.CreateRepository(ctx, docker.RepositoryInput{
		NamespaceID: namespaceId,
		Name:        name,
	}, docker.RepositoryDetailsInput{
		ShortDescription: shortDescription,
		FullDescription:  fullDescription,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("createRepository: %w", err))
	}
	d.SetId(repo.ID)
	return resourceDockerRepositoryRead(ctx, d, m)
}
