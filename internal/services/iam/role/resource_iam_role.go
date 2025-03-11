package role

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var descriptions = map[string]string{
	"role": "Organization administrators can create roles for the users, services, and devices that interact with their organizations. Roles are a collection of permissions and provide a way to manage the assignment of permissions to users. Permissions are privileges that define what a user is allowed to do. The roles can contain permissions from one or more applications and services. That is, if there are two products being used, permissions from both products can be added to one role. Roles are assigned to a group",
}

func ResourceIAMRole() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["role"],
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		CreateContext: resourceIAMRoleCreate,
		ReadContext:   resourceIAMRoleRead,
		UpdateContext: resourceIAMRoleUpdate,
		DeleteContext: resourceIAMRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tools.ValidateUpperString,
				Description:  "The role name.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
				Description:      "The role description.",
			},
			"managing_organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managing organization of the role.",
			},
			"permissions": {
				Type:        schema.TypeSet,
				MaxItems:    100,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of permissions IDs assigned to this role.",
			},
			"ticket_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Removal protection of some ticket only permissions.",
			},
		},
	}
}

func resourceIAMRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	managingOrganization := d.Get("managing_organization").(string)
	permissions := tools.ExpandStringList(d.Get("permissions").(*schema.Set).List())

	var role *iam.Role
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		role, resp, err = client.Roles.CreateRole(name, description, managingOrganization)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		if resp.StatusCode() != http.StatusConflict {
			return diag.FromErr(err)
		}
		// Already exists most likely, adopt it
		var roles *[]iam.Role
		roles, _, err = client.Roles.GetRoles(&iam.GetRolesOptions{
			Name:           &name,
			OrganizationID: &managingOrganization,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(*roles) == 0 || (*roles)[0].ManagingOrganization != managingOrganization {
			return diag.FromErr(fmt.Errorf("conflict creating role, mismatched managing_organization: '%s' != '%s'",
				(*roles)[0].ManagingOrganization, managingOrganization))
		}
		role = &(*roles)[0]
	}
	for _, p := range permissions {
		result, resp, err := client.Roles.AddRolePermission(*role, p)
		if err != nil {
			// Clean up
			_, _, _ = client.Roles.DeleteRole(*role)
			return diag.FromErr(fmt.Errorf("error adding permission '%s': %v", p, err))
		}
		if resp == nil || !(resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusMultiStatus) {
			// Clean up
			_, _, _ = client.Roles.DeleteRole(*role)
			return diag.FromErr(fmt.Errorf("error adding permission '%s': %v", p, result))
		}
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "invalid permission",
				Detail:   fmt.Sprintf("permission '%s' is invalid", p),
			})
		}
	}
	d.SetId(role.ID)
	res := resourceIAMRoleRead(ctx, d, m)
	diags = append(diags, res...)
	return diags
}

func resourceIAMRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	role, resp, err := client.Roles.GetRoleByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		if resp != nil && resp.StatusCode() == http.StatusForbidden { // See INC0080073
			orgId := d.Get("managing_organization").(string)
			if client.HasPermissions(orgId, "ROLE.WRITE") {
				// If we have ROLE.WRITE permission and get HTTP 403 we conclude the Role is gone
				d.SetId("")
				return diags
			}
		}
		return diag.FromErr(err)
	}
	_ = d.Set("description", role.Description)
	_ = d.Set("name", role.Name)
	_ = d.Set("managing_organization", role.ManagingOrganization)

	var permissions *[]string

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		permissions, resp, err = client.Roles.GetRolePermissions(*role)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp.StatusCode() == http.StatusForbidden { // IAM limitation
			return diags // Use Terraform as source of truth
		}
		return diag.FromErr(err)
	}
	_ = d.Set("permissions", permissions)
	return diags
}

func resourceIAMRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	role, _, err := client.Roles.GetRoleByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("description") {
		return diag.FromErr(fmt.Errorf("description changes are not supported"))
	}

	if d.HasChange("permissions") {
		o, n := d.GetChange("permissions")
		oldList := tools.ExpandStringList(o.(*schema.Set).List())
		newList := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newList, oldList)
		toRemove := tools.Difference(oldList, newList)

		res := addAndRemovePermissions(ctx, *role, toAdd, toRemove, client)
		diags = append(diags, res...)
	}
	return diags
}

func addAndRemovePermissions(_ context.Context, role iam.Role, toAdd, toRemove []string, client *iam.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	// Additions
	if len(toAdd) > 0 {
		for _, v := range toAdd {
			_, resp, err := client.Roles.AddRolePermission(role, v)
			if err != nil {
				if resp != nil && resp.StatusCode() == http.StatusNotFound {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "invalid permission",
						Detail:   fmt.Sprintf("permission '%s' is invalid", v),
					})
					continue
				}
				return diag.FromErr(err)
			}
		}
	}
	// Removals
	for _, v := range toRemove {
		_, resp, err := client.Roles.RemoveRolePermission(role, v)
		if err != nil {
			if resp != nil && resp.StatusCode() == http.StatusNotFound {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "invalid permission",
					Detail:   fmt.Sprintf("permission '%s' is invalid", v),
				})
				continue // Accept 404 in case of invalid permissions
			}
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceIAMRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var role iam.Role
	role.ID = d.Id()

	var resp *iam.Response
	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		_, resp, err = client.Roles.DeleteRole(role)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil || resp.StatusCode() != http.StatusNoContent {
		return diag.FromErr(config.ErrDeleteRoleFailed)
	}
	d.SetId("")
	return diags
}
