package namespace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceDockerNamespaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDockerNamespacesRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_repos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}

}

func dataSourceDockerNamespacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespaces, err := client.Namespaces.GetNamespaces(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("docker_namespaces")

	var names []string
	var numRepos []int

	for _, namespace := range *namespaces {
		names = append(names, namespace.ID)
		numRepos = append(numRepos, namespace.NumRepos)
	}
	_ = d.Set("names", names)
	_ = d.Set("num_repos", numRepos)

	return diags
}
