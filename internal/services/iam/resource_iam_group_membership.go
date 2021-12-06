package iam

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceIAMGroupMembership() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceIAMGroupMembershipCreate,
		ReadContext:   resourceIAMGroupMembershipRead,
		UpdateContext: resourceIAMGroupMembershipUpdate,
		DeleteContext: resourceIAMGroupMembershipDelete,

		Schema: map[string]*schema.Schema{
			"iam_group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
			},
			"users": {
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
			"services": {
				Type:     schema.TypeSet,
				MaxItems: 2000,
				Optional: true,
				Elem:     tools.StringSchema(),
			},
		},
	}
}

func resourceIAMGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	groupId := d.Get("iam_group_id").(string)

	group, resp, err := client.Groups.GetGroupByID(groupId)

	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusForbidden:
				err = fmt.Errorf("no permission to read group details: %w", err)
			default:
				err = fmt.Errorf("error reading group '%s' (HTTP %d): %w", groupId, resp.StatusCode, err)
			}
		}
		return diag.FromErr(err)
	}
	// Add users
	users := tools.ExpandStringList(d.Get("users").(*schema.Set).List())
	if len(users) > 0 {
		err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
			result, resp, err := client.Groups.AddMembers(*group, users...)
			if err != nil {
				_ = client.TokenRefresh()
			}
			if resp == nil {
				return nil, err
			}
			if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMultiStatus) {
				return resp.Response, backoff.Permanent(fmt.Errorf("failed to add members: %v %w", result, err))
			}
			return resp.Response, err
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding users: %w", err))
		}
	}

	// Add services
	services := tools.ExpandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
			result, resp, err := client.Groups.AddServices(*group, services...)
			if err != nil {
				_ = client.TokenRefresh()
			}
			if resp == nil {
				return nil, err
			}
			if err != nil {
				return resp.Response, err
			}
			if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMultiStatus) {
				return resp.Response, backoff.Permanent(fmt.Errorf("failed to add services: %v", result))
			}
			return resp.Response, err
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding services: %v", err))
		}
	}
	d.SetId(uuid.NewString())
	return diags
}

func resourceIAMGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	groupId := d.Get("iam_group_id").(string)

	group, resp, err := client.Groups.GetGroupByID(groupId)

	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusForbidden:
				err = fmt.Errorf("no permission to read group details: %w", err)
			default:
				err = fmt.Errorf("error reading group '%s' (HTTP %d): %w", groupId, resp.StatusCode, err)
			}
		}
		return diag.FromErr(err)
	}
	// Remove users
	users := tools.ExpandStringList(d.Get("users").(*schema.Set).List())
	if len(users) > 0 {
		err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
			result, resp, err := client.Groups.RemoveMembers(*group, users...)
			if resp == nil {
				return nil, err
			}
			if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMultiStatus) {
				return resp.Response, backoff.Permanent(fmt.Errorf("failed to remove members: %v %w", result, err))
			}
			return resp.Response, err
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error removing users: %w", err))
		}
	}

	// Remove services
	services := tools.ExpandStringList(d.Get("services").(*schema.Set).List())
	if len(services) > 0 {
		err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
			result, resp, err := client.Groups.RemoveServices(*group, services...)
			if resp == nil {
				return nil, err
			}
			if err != nil {
				return resp.Response, err
			}
			if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMultiStatus) {
				return resp.Response, backoff.Permanent(fmt.Errorf("failed to remove services: %v", result))
			}
			return resp.Response, err
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error removing services: %v", err))
		}
	}
	d.SetId("")
	return diags
}

func resourceIAMGroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	groupId := d.Get("iam_group_id").(string)

	var group iam.Group
	group.ID = groupId

	// Users
	if d.HasChange("users") {
		o, n := d.GetChange("users")
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newList := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newList, old)
		toRemove := tools.Difference(old, newList)

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
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newList := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newList, old)
		toRemove := tools.Difference(old, newList)

		if len(toRemove) > 0 {
			err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
				_, resp, err := client.Groups.RemoveServices(group, toRemove...)
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
			err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
				_, resp, err := client.Groups.AddServices(group, toAdd...)
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
	return diags
}

func resourceIAMGroupMembershipRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	groupId := d.Get("iam_group_id").(string)

	group, resp, err := client.Groups.GetGroupByID(groupId)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	// Users
	// We only deal with users we know
	users, _, err := client.Users.GetUsers(&iam.GetUserOptions{
		GroupID: &group.ID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving users from group: %v", err))
	}
	var presentUsers []string
	knownUsers := tools.ExpandStringList(d.Get("users").(*schema.Set).List())
	for _, u := range knownUsers {
		if tools.ContainsString(users.UserUUIDs, u) {
			presentUsers = append(presentUsers, u)
		}
	}
	_ = d.Set("users", tools.SchemaSetStrings(presentUsers))

	// Services
	// We only deal with services we know
	var verifiedServices []string
	services := tools.ExpandStringList(d.Get("services").(*schema.Set).List())
	for _, service := range services {
		groups, _, err := client.Groups.GetGroups(&iam.GetGroupOptions{
			MemberType: tools.String("SERVICE"),
			MemberID:   &service,
		})
		if err != nil || groups == nil {
			continue
		}
		for _, g := range *groups {
			if g.ID == group.ID {
				verifiedServices = append(verifiedServices, service)
				continue
			}
		}
	}
	_ = d.Set("services", tools.SchemaSetStrings(verifiedServices))
	return diags
}
