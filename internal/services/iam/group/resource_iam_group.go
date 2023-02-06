package group

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

var descriptions = map[string]string{
	"group": "Organization administrators can create groups to help organize users by teams, projects, or user function. Groups must be uniquely named within the organization in which they were created. Admins can add a single identity to one or more groups. They can also assign one or more roles to a group. All users in a group will share the permissions that are associated with the role.",
}

func ResourceIAMGroup() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["group"],
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 3,
		CreateContext: resourceIAMGroupCreate,
		ReadContext:   resourceIAMGroupRead,
		UpdateContext: resourceIAMGroupUpdate,
		DeleteContext: resourceIAMGroupDelete,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    ResourceIAMGroupV0().CoreConfigSchema().ImpliedType(),
				Upgrade: patchIAMGroupV0,
				Version: 0,
			},
			{
				Type:    ResourceIAMGroupV1().CoreConfigSchema().ImpliedType(),
				Upgrade: patchIAMGroupV1,
				Version: 1,
			},
			{
				Type:    ResourceIAMGroupV2().CoreConfigSchema().ImpliedType(),
				Upgrade: patchIAMGroupV2,
				Version: 2,
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
				Description:      "The group name.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
				Description:      "The group description.",
			},
			"managing_organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managing organization ID.",
			},
			"roles": {
				Type:        schema.TypeSet,
				MaxItems:    1000,
				Required:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of role IDS to assign to this group.",
			},
			"users": {
				Type:        schema.TypeSet,
				MaxItems:    2000,
				Optional:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of user IDs to include in this group. The provider only manages this list of users. Existing users added by others means to the group by the provider. It is not practical to manage hundreds or thousands of users this way of course.",
			},
			"services": {
				Type:        schema.TypeSet,
				MaxItems:    2000,
				Optional:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of service identity IDs to include in this group.",
			},
			"devices": {
				Type:        schema.TypeSet,
				MaxItems:    2000,
				Optional:    true,
				Elem:        tools.StringSchema(),
				Description: "The list of IAM device identity IDs to include in this group.",
			},
			"drift_detection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When enabled, the provider will perform additional API calls to determine if any changes were made outside of Terraform to user and service assignments of this group.",
			},
			"iam_device_bug_workaround": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Deprecated:  "This workaround is no longer required and will be removed in the near future.",
				Description: "Deprecated, do not use.",
			},
		},
	}
}

func resourceIAMGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var group iam.Group
	group.Description = d.Get("description").(string)
	group.Name = d.Get("name").(string)
	group.ManagingOrganization = d.Get("managing_organization").(string)

	var createdGroup *iam.Group
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var resp *iam.Response
		var err error
		createdGroup, resp, err = client.Groups.CreateGroup(group)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})

	if err != nil {
		return diag.FromErr(err)
	}
	roles := tools.ExpandStringList(d.Get("roles").(*schema.Set).List())

	_ = d.Set("name", createdGroup.Name)
	_ = d.Set("description", createdGroup.Description)
	_ = d.Set("managing_organization", createdGroup.ManagingOrganization)

	// Add roles
	for _, r := range roles {
		role, _, _ := client.Roles.GetRoleByID(r)
		if role != nil {
			err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
				_, resp, err := client.Groups.AssignRole(ctx, *createdGroup, *role)
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			}, append(tools.StandardRetryOnCodes, http.StatusUnprocessableEntity)...) // Handle intermittent HTTP 422 errors
			if err != nil {
				// Cleanup
				_ = purgeGroupContent(ctx, client, createdGroup.ID, d)
				_, _, _ = client.Groups.DeleteGroup(*createdGroup)
				return diag.FromErr(fmt.Errorf("error adding roles: %v", err))
			}
		}
	}

	// Add users
	users := tools.ExpandStringList(d.Get("users").(*schema.Set).List())
	if len(users) > 0 {
		err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
			result, resp, err := client.Groups.AddMembers(ctx, *createdGroup, users...)
			if resp == nil {
				return nil, err
			}
			if !(resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusMultiStatus) {
				return resp.Response, backoff.Permanent(fmt.Errorf("failed to add members: %v %w", result, err))
			}
			return resp.Response, err
		})
		if err != nil {
			// Cleanup
			_ = purgeGroupContent(ctx, client, createdGroup.ID, d)
			_, _, _ = client.Groups.DeleteGroup(*createdGroup)
			return diag.FromErr(fmt.Errorf("error adding users: %w", err))
		}
	}

	// Add services
	services := tools.ExpandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
			result, resp, err := client.Groups.AddServices(ctx, *createdGroup, services...)
			if resp == nil {
				return nil, err
			}
			if err != nil {
				return resp.Response, err
			}
			if !(resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusMultiStatus) {
				return resp.Response, backoff.Permanent(fmt.Errorf("failed to add services: %v", result))
			}
			return resp.Response, err
		})
		if err != nil {
			// Cleanup
			_ = purgeGroupContent(ctx, client, createdGroup.ID, d)
			_, _, _ = client.Groups.DeleteGroup(*createdGroup)
			return diag.FromErr(fmt.Errorf("error adding services: %v", err))
		}
	}

	// Add devices
	devices := tools.ExpandStringList(d.Get("devices").(*schema.Set).List())
	if len(devices) > 0 {
		err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
			result, resp, err := client.Groups.AddDevices(ctx, *createdGroup, devices...)
			if resp == nil {
				return nil, err
			}
			if err != nil {
				return resp.Response, err
			}
			if !(resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusMultiStatus) {
				return resp.Response, backoff.Permanent(fmt.Errorf("failed to add devices: %v", result))
			}
			return resp.Response, err
		})
		if err != nil {
			// Cleanup
			_ = purgeGroupContent(ctx, client, createdGroup.ID, d)
			_, _, _ = client.Groups.DeleteGroup(*createdGroup)
			return diag.FromErr(fmt.Errorf("error adding devices: %v", err))
		}
	}

	d.SetId(createdGroup.ID)
	return resourceIAMGroupRead(ctx, d, m)
}

func resourceIAMGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	driftDetection := d.Get("drift_detection").(bool)

	group, resp, err := client.Groups.GetGroupByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
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
	var roleIDs []string
	for _, r := range *roles {
		roleIDs = append(roleIDs, r.ID)
	}
	_ = d.Set("roles", tools.SchemaSetStrings(roleIDs))

	if driftDetection { // Only do drift detection when explicitly enabled
		// Users
		users, err := getGroupResourcesByMemberType(ctx, client, group.ID, iam.GroupMemberTypeUser)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving users from group: %v", err))
		}
		_ = d.Set("users", tools.SchemaSetStrings(users))

		// Services
		services, err := getGroupResourcesByMemberType(ctx, client, group.ID, iam.GroupMemberTypeService)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving services from group: %v", err))
		}
		_ = d.Set("services", tools.SchemaSetStrings(services))

		// Devices
		devices, err := getGroupResourcesByMemberType(ctx, client, group.ID, iam.GroupMemberTypeDevice)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving devices from group: %v", err))
		}
		_ = d.Set("devices", tools.SchemaSetStrings(devices))
	}
	return diags
}

func resourceIAMGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
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
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newList := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newList, old)
		toRemove := tools.Difference(old, newList)

		if len(toRemove) > 0 {
			_, _, _ = client.Groups.RemoveMembers(ctx, group, toRemove...)
		}
		if len(toAdd) > 0 {
			_, _, _ = client.Groups.AddMembers(ctx, group, toAdd...)
		}
	}

	// Services
	if d.HasChange("services") {
		o, n := d.GetChange("services")
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newList := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newList, old)
		toRemove := tools.Difference(old, newList)

		if len(toRemove) > 0 {
			err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
				_, resp, err := client.Groups.RemoveServices(ctx, group, toRemove...)
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		if len(toAdd) > 0 {
			err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
				_, resp, err := client.Groups.AddServices(ctx, group, toAdd...)
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}

	// Devices
	if d.HasChange("devices") {
		o, n := d.GetChange("devices")
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newList := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newList, old)
		toRemove := tools.Difference(old, newList)

		if len(toRemove) > 0 {
			err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
				_, resp, err := client.Groups.RemoveDevices(ctx, group, toRemove...)
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		if len(toAdd) > 0 {
			err = tools.TryHTTPCall(ctx, 5, func() (*http.Response, error) {
				_, resp, err := client.Groups.AddDevices(ctx, group, toAdd...)
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
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
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newValues := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newValues, old)
		toRemove := tools.Difference(old, newValues)

		// Handle additions
		if len(toAdd) > 0 {
			for _, v := range toAdd {
				var role = iam.Role{ID: v}
				_, _, err := client.Groups.AssignRole(ctx, group, role)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// Remove every role. Simpler to remove and add newValues ones,
		for _, v := range toRemove {
			var role = iam.Role{ID: v}
			_, resp, err := client.Groups.RemoveRole(ctx, group, role)
			if resp != nil && resp.StatusCode() == http.StatusUnprocessableEntity {
				return nil // Role is already gone
			}
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	return diags
}

func purgeGroupContent(ctx context.Context, client *iam.Client, id string, d *schema.ResourceData) error {
	var group iam.Group
	group.ID = id

	// Remove all users first before attempting delete
	users, _, err := client.Users.GetAllUsers(&iam.GetUserOptions{
		GroupID: &group.ID,
	})
	if err != nil {
		return fmt.Errorf("retrieving user list of group %s: %w", group.ID, err)
	}
	if len(users) > 0 {
		for _, u := range users {
			_ = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
				_, resp, err := client.Groups.RemoveMembers(ctx, group, u)
				if resp != nil && resp.StatusCode() == http.StatusUnprocessableEntity {
					return resp.Response, nil // User is already gone
				}
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			})
		}
	}

	// Remove all services first before attempting delete
	services := tools.ExpandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		for _, s := range services {
			_ = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
				_, resp, err := client.Groups.RemoveServices(ctx, group, s)
				if resp != nil && resp.StatusCode() == http.StatusUnprocessableEntity {
					return resp.Response, nil // Service is already gone
				}
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			})
		}
	}

	// Remove all devices first before attempting delete
	devices := tools.ExpandStringList(d.Get("devices").(*schema.Set).List())
	if len(devices) > 0 {
		for _, s := range devices {
			_ = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
				_, resp, err := client.Groups.RemoveDevices(ctx, group, s)
				if resp != nil && resp.StatusCode() == http.StatusUnprocessableEntity {
					return resp.Response, nil // Service is already gone
				}
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			})
		}
	}

	// Remove all associated roles
	roles, _, err := client.Roles.GetRoles(&iam.GetRolesOptions{
		GroupID: &group.ID,
	})
	if err != nil || roles == nil {
		return fmt.Errorf("retrieving roles of group %s: %w", group.ID, err)
	}
	if len(*roles) > 0 {
		for _, r := range *roles {
			_ = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
				var role = iam.Role{ID: r.ID}
				_, resp, err := client.Groups.RemoveRole(ctx, group, role)
				if resp != nil && resp.StatusCode() == http.StatusUnprocessableEntity {
					return resp.Response, nil // Role is already gone
				}
				if resp == nil {
					return nil, err
				}
				return resp.Response, err
			})
		}
	}
	return nil
}

func resourceIAMGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var group iam.Group
	group.ID = d.Id()
	if err := purgeGroupContent(ctx, client, group.ID, d); err != nil {
		return diag.FromErr(fmt.Errorf("error purging group content: %v", err))
	}

	// Query group to sync it up (to force IAM sync?)
	_ = resourceIAMGroupRead(ctx, d, m)

	var ok bool
	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var resp *iam.Response
		var err error
		ok, resp, err = client.Groups.DeleteGroup(group)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrDeleteGroupFailed)
	}
	d.SetId("")
	return diags
}
