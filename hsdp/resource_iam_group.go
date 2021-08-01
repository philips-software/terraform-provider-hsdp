package hsdp

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMGroup() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMGroupCreate,
		ReadContext:   resourceIAMGroupRead,
		UpdateContext: resourceIAMGroupUpdate,
		DeleteContext: resourceIAMGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressCaseDiffs,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressWhenGenerated,
			},
			"managing_organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIAMGroupCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var group iam.Group
	group.Description = d.Get("description").(string)
	group.Name = d.Get("name").(string)
	group.ManagingOrganization = d.Get("managing_organization").(string)

	var createdGroup *iam.Group
	err = tryIAMCall(func() (*iam.Response, error) {
		var resp *iam.Response
		var err error
		createdGroup, resp, err = client.Groups.CreateGroup(group)
		return resp, err
	})

	if err != nil {
		return diag.FromErr(err)
	}
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	d.SetId(createdGroup.ID)
	_ = d.Set("name", createdGroup.Name)
	_ = d.Set("description", createdGroup.Description)
	_ = d.Set("managing_organization", createdGroup.ManagingOrganization)

	// Add roles
	for _, r := range roles {
		role, _, _ := client.Roles.GetRoleByID(r)
		if role != nil {
			err = tryIAMCall(func() (*iam.Response, error) {
				_, resp, err := client.Groups.AssignRole(*createdGroup, *role)
				return resp, err
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}

	// Add users
	users := expandStringList(d.Get("users").(*schema.Set).List())
	if len(users) > 0 {
		err = tryIAMCall(func() (*iam.Response, error) {
			_, resp, err := client.Groups.AddMembers(*createdGroup, users...)
			return resp, err
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Add services
	services := expandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		err = tryIAMCall(func() (*iam.Response, error) {
			_, resp, err := client.Groups.AddServices(*createdGroup, services...)
			return resp, err
		})
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}

func resourceIAMGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	group, resp, err := client.Groups.GetGroupByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("managing_organization", group.ManagingOrganization)
	_ = d.Set("description", group.Description)
	_ = d.Set("name", group.Name)
	roles, _, err := client.Groups.GetRoles(*group)
	if err != nil {
		return diag.FromErr(err)
	}
	roleIDs := make([]string, len(*roles))
	for i, r := range *roles {
		roleIDs[i] = r.ID
	}
	_ = d.Set("roles", &roleIDs)
	return diags
}

func resourceIAMGroupUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var group iam.Group
	group.ID = d.Id()
	if d.HasChange("description") {
		group.Description = d.Get("description").(string)
		_, _, err := client.Groups.UpdateGroup(group)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// Users
	if d.HasChange("users") {
		o, n := d.GetChange("users")
		old := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, old)
		toRemove := difference(old, newList)

		if len(toRemove) > 0 {
			_, _, _ = client.Groups.RemoveMembers(group, toRemove...)
		}
		if len(toAdd) > 0 {
			_, _, _ = client.Groups.AddMembers(group, toAdd...)
		}
	}

	// Services
	if d.HasChange("services") {
		o, n := d.GetChange("services")
		old := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, old)
		toRemove := difference(old, newList)

		if len(toRemove) > 0 {
			err = tryIAMCall(func() (*iam.Response, error) {
				_, resp, err := client.Groups.RemoveServices(group, toRemove...)
				return resp, err
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		if len(toAdd) > 0 {
			err = tryIAMCall(func() (*iam.Response, error) {
				_, resp, err := client.Groups.AddServices(group, toAdd...)
				return resp, err
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}
	if len(diags) > 0 {
		return diags
	}

	// Roles
	if d.HasChange("roles") {
		o, n := d.GetChange("roles")
		old := expandStringList(o.(*schema.Set).List())
		newValues := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newValues, old)
		toRemove := difference(old, newValues)

		// Handle additions
		if len(toAdd) > 0 {
			for _, v := range toAdd {
				var role = iam.Role{ID: v}
				_, _, err := client.Groups.AssignRole(group, role)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// Remove every role. Simpler to remove and add newValues ones,
		for _, v := range toRemove {
			var role = iam.Role{ID: v}
			_, _, err := client.Groups.RemoveRole(group, role)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	return diags
}

func resourceIAMGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var group iam.Group
	group.ID = d.Id()

	// Remove all users first before attempting delete
	users := expandStringList(d.Get("users").(*schema.Set).List())
	if len(users) > 0 {
		for _, u := range users {
			err := tryIAMCall(func() (*iam.Response, error) {
				_, resp, err := client.Groups.RemoveMembers(group, u)
				if resp != nil && resp.StatusCode == http.StatusUnprocessableEntity {
					return resp, nil // User is already gone
				}
				return resp, err
			}, http.StatusInternalServerError)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		if len(diags) > 0 {
			return diags
		}
	}

	// Remove all services first before attempting delete
	services := expandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		for _, s := range services {
			err := tryIAMCall(func() (*iam.Response, error) {
				_, resp, err := client.Groups.RemoveServices(group, s)
				if resp != nil && resp.StatusCode == http.StatusUnprocessableEntity {
					return resp, nil // Service is already gone
				}
				return resp, err
			}, http.StatusInternalServerError)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		if len(diags) > 0 {
			return diags
		}
	}

	// Remove all associated roles
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	if len(roles) > 0 {
		for _, r := range roles {
			err := tryIAMCall(func() (*iam.Response, error) {
				var role = iam.Role{ID: r}
				_, resp, err := client.Groups.RemoveRole(group, role)
				if resp != nil && resp.StatusCode == http.StatusUnprocessableEntity {
					return resp, nil // Role is already gone
				}
				return resp, err
			}, http.StatusInternalServerError)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		if len(diags) > 0 {
			return diags
		}
	}

	// Query group to sync it up (to force IAM sync?)
	_ = resourceIAMGroupRead(ctx, d, m)

	var ok bool
	err = tryIAMCall(func() (*iam.Response, error) {
		var resp *iam.Response
		var err error
		ok, resp, err = client.Groups.DeleteGroup(group)
		return resp, err
	}, http.StatusInternalServerError)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(ErrDeleteGroupFailed)
	}
	d.SetId("")
	return diags
}
