package namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/console/docker"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceDockerNamespace() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDockerNamespaceCreate,
		ReadContext:   resourceDockerNamespaceRead,
		DeleteContext: resourceDockerNamespaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_repos": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDockerNamespaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	err = client.Namespaces.DeleteNamespace(ctx, docker.Namespace{ID: id})
	if err != nil {
		return diag.FromErr(fmt.Errorf("deleteNamespace: %w", err))
	}
	d.SetId("")
	return diags
}

func resourceDockerNamespaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()

	ns, err := client.Namespaces.GetNamespaceByID(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("getNamespaceById: %w", err))
	}
	_ = d.Set("created_at", ns.CreatedAt.Format(time.RFC3339))
	_ = d.Set("num_repos", ns.NumRepos)
	_ = d.Set("is_public", ns.IsPublic)

	return diags
}

func resourceDockerNamespaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	ns, err := client.Namespaces.CreateNamespace(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("createNamespace: %w", err))
	}
	d.SetId(ns.ID)
	return resourceDockerNamespaceRead(ctx, d, m)
}
