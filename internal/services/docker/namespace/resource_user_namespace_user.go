package namespace

import (
	"context"
	"fmt"

	"github.com/dip-software/go-dip-api/console/docker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceDockerNamespaceUser() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDockerNamespaceUserCreate,
		ReadContext:   resourceDockerNamespaceUserRead,
		UpdateContext: resourceDockerNamespaceUserUpdate,
		DeleteContext: resourceDockerNamespaceUserDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"can_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_pull": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"can_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDockerNamespaceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var id int
	_, _ = fmt.Sscanf(d.Id(), "%d", &id)

	canDelete := d.Get("can_delete").(bool)
	canPush := d.Get("can_push").(bool)
	canPull := d.Get("can_pull").(bool)
	isAdmin := d.Get("is_admin").(bool)

	err = client.Namespaces.UpdateNamespaceUserAccess(ctx, id, docker.UserNamespaceAccessInput{
		CanDelete: canDelete,
		CanPush:   canPush,
		CanPull:   canPull,
		IsAdmin:   isAdmin,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("updateNamespaceUserAccess: %w", err))
	}
	return resourceDockerNamespaceUserRead(ctx, d, m)
}

func resourceDockerNamespaceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	userId := d.Get("user_id").(string)
	namespaceId := d.Get("namespace_id").(string)

	err = client.Namespaces.DeleteNamespaceUser(ctx, namespaceId, userId)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deleteNamespaceUser: %w", err))
	}
	d.SetId("")
	return diags
}

func resourceDockerNamespaceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespaceId := d.Get("namespace_id").(string)
	username := d.Get("username").(string)

	users, err := client.Namespaces.GetNamespaceUsers(ctx, docker.Namespace{ID: namespaceId})
	if err != nil {
		return diag.FromErr(fmt.Errorf("dockerNamespaceUserRead(id): %w", err))
	}
	var namespaceUser *docker.NamespaceUser
	for _, user := range *users {
		if user.Username == username {
			namespaceUser = &user
			break
		}
	}
	if namespaceUser == nil {
		d.SetId("")
		return diags
	}
	_ = d.Set("can_push", namespaceUser.NamespaceAccess.CanPush)
	_ = d.Set("can_delete", namespaceUser.NamespaceAccess.CanDelete)
	_ = d.Set("can_pull", namespaceUser.NamespaceAccess.CanPull)
	_ = d.Set("is_admin", namespaceUser.NamespaceAccess.IsAdmin)

	return diags
}

func resourceDockerNamespaceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var client *docker.Client
	var err error

	client, err = c.DockerClient()
	if err != nil {
		return diag.FromErr(err)
	}
	username := d.Get("username").(string)
	namespaceId := d.Get("namespace_id").(string)
	canDelete := d.Get("can_delete").(bool)
	canPush := d.Get("can_push").(bool)
	canPull := d.Get("can_pull").(bool)
	isAdmin := d.Get("is_admin").(bool)

	ns, err := client.Namespaces.AddNamespaceUser(ctx, namespaceId, username, docker.UserNamespaceAccessInput{
		CanDelete: canDelete,
		CanPush:   canPush,
		CanPull:   canPull,
		IsAdmin:   isAdmin,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("createNamespace: %w", err))
	}
	_ = d.Set("user_id", ns.UserID)
	d.SetId(fmt.Sprintf("%d", ns.ID))
	return resourceDockerNamespaceUserRead(ctx, d, m)
}
