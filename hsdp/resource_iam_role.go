package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMRole() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMRoleCreate,
		ReadContext:   resourceIAMRoleRead,
		UpdateContext: resourceIAMRoleUpdate,
		DeleteContext: resourceIAMRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUpperString,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ticket_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	managingOrganization := d.Get("managing_organization").(string)
	permissions := expandStringList(d.Get("permissions").(*schema.Set).List())

	role, resp, err := client.Roles.CreateRole(name, description, managingOrganization)
	if err != nil {
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		if resp.StatusCode != http.StatusConflict {
			return diag.FromErr(err)
		}
		// Already exists most likely, adopt it
		var roles *[]iam.Role
		roles, _, err = client.Roles.GetRoles(&iam.GetRolesOptions{
			Name: &name,
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
		_, _, _ = client.Roles.AddRolePermission(*role, p)
	}
	d.SetId(role.ID)
	readDiags := resourceIAMRoleRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}
	return diags
}

func resourceIAMRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	role, resp, err := client.Roles.GetRoleByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("description", role.Description)
	_ = d.Set("name", role.Name)
	_ = d.Set("managing_organization", role.ManagingOrganization)

	permissions, _, err := client.Roles.GetRolePermissions(*role)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("permissions", permissions)
	d.SetId(role.ID)
	return diags
}

func resourceIAMRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
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
		oldList := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, oldList)
		toRemove := difference(oldList, newList)

		// Additions
		if len(toAdd) > 0 {
			for _, v := range toAdd {
				_, _, err := client.Roles.AddRolePermission(*role, v)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// Removals
		for _, v := range toRemove {
			ticketProtection := d.Get("ticket_protection").(bool)
			if ticketProtection && v == "CLIENT.SCOPES" {
				return diag.FromErr(fmt.Errorf("Refusing to remove CLIENT.SCOPES permission, set ticket_protection to `false` to override"))
			}
			_, _, err := client.Roles.RemoveRolePermission(*role, v)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	return diags
}

func resourceIAMRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var role iam.Role
	role.ID = d.Id()

	_, _, _ = client.Roles.DeleteRole(role)
	// Note: 2020-12-17, the DeleteRole call will always fail. Thanks IAM.
	d.SetId("")
	return diags
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(int)
		if ok && val != 0 {
			vs = append(vs, val)
		}
	}
	return vs
}
