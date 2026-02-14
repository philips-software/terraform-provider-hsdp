package service_key

import (
	"context"
	"fmt"
	"time"

	"github.com/philips-software/go-dip-api/console/docker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceDockerServiceKey() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDockerServiceKeyCreate,
		ReadContext:   resourceDockerServiceKeyRead,
		DeleteContext: resourceDockerServiceKeyDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDockerServiceKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var resourceID int
	_, _ = fmt.Sscanf(d.Id(), "%d", &resourceID)
	err = client.ServiceKeys.DeleteServiceKey(ctx, docker.ServiceKey{ID: resourceID})
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete service key error: %w", err))
	}
	d.SetId("")
	return diags
}

func resourceDockerServiceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var resourceID int
	_, _ = fmt.Sscanf(d.Id(), "%d", &resourceID)

	key, err := client.ServiceKeys.GetServiceKeyByID(ctx, resourceID)
	if err != nil {
		// Assume service key was deleted
		d.SetId("")
		return diags
	}
	_ = d.Set("description", key.Description)
	_ = d.Set("username", key.Username)
	_ = d.Set("created_at", key.CreatedAt.Format(time.RFC3339))
	return diags
}

func resourceDockerServiceKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	description := d.Get("description").(string)
	created, err := client.ServiceKeys.CreateServiceKey(ctx, description)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("password", created.Password) // Only time we get the password
	d.SetId(fmt.Sprintf("%d", created.ID))
	return resourceDockerServiceKeyRead(ctx, d, m)
}
