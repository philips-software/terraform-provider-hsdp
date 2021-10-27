package namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceDockerNamespace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDockerNamespaceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"num_repos": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceDockerNamespaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	ns, err := client.Namespaces.GetNamespaceByID(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)

	_ = d.Set("created_at", ns.CreatedAt.Format(time.RFC3339))
	_ = d.Set("num_repos", ns.NumRepos)
	_ = d.Set("is_public", ns.IsPublic)
	repos, err := client.Namespaces.GetRepositories(ctx, ns.ID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading repositories: %w", err))
	}
	var repositories []string
	for _, r := range *repos {
		repositories = append(repositories, r.Name)
	}
	_ = d.Set("repositories", repositories)

	return diags
}
