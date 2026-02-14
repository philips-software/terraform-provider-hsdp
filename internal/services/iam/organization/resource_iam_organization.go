package organization

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/go-dip-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var descriptions = map[string]string{
	"organization": "An organization is a container that owns sub-organizations, groups, and users, as well as other identities like devices and services.",
}

func ResourceIAMOrg() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["organization"],
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 4,
		CreateContext: resourceIAMOrgCreate,
		ReadContext:   resourceIAMOrgRead,
		UpdateContext: resourceIAMOrgUpdate,
		DeleteContext: resourceIAMOrgDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the IAM Organization.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the organization.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name to use for this organization.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The organization type.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional external identifier for this organization.",
			},
			"is_root_org": {
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "This field is deprecated, please remove it",
				Description: "Deprecated, do not use.",
			},
			"parent_org_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent organization ID.",
			},
			"wait_for_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Blocks until the organization delete has completed. Default: false. The organization delete process can take some time as all its associated resources like users, groups, roles etc. are removed recursively. This option is useful for ephemeral environments where the same organization might be recreated shortly after a destroy operation.",
			},
			"active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Weather the organization is active or not.",
			},
		},
	}
}

func resourceIAMOrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	externalID := d.Get("external_id").(string)
	orgType := d.Get("type").(string)
	displayName := d.Get("display_name").(string)
	parentOrgID, ok := d.Get("parent_org_id").(string)
	if !ok {
		return diag.FromErr(config.ErrMissingParentOrgID)
	}
	var newOrg iam.Organization
	newOrg.Name = name
	newOrg.Description = description
	newOrg.Parent.Value = parentOrgID
	newOrg.ExternalID = externalID
	newOrg.DisplayName = displayName
	newOrg.Type = orgType

	var org *iam.Organization
	var resp *iam.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		org, resp, err = client.Organizations.CreateOrganization(newOrg)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if org == nil {
		return diag.FromErr(fmt.Errorf("failed to create organization: %d", resp.StatusCode()))
	}
	d.SetId(org.ID)
	return resourceIAMOrgRead(ctx, d, m)
}

func resourceIAMOrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	var resp *iam.Response
	var org *iam.Organization

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		org, resp, err = client.Organizations.GetOrganizationByID(id)
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound { // Gone
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	_ = d.Set("description", org.Description)
	_ = d.Set("name", org.Name)
	_ = d.Set("external_id", org.ExternalID)
	_ = d.Set("parent_org_id", org.Parent.Value)
	_ = d.Set("display_name", org.DisplayName)
	_ = d.Set("active", org.Active)
	_ = d.Set("type", org.Type)
	return diags
}

func resourceIAMOrgUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("description") || d.HasChange("name") ||
		d.HasChange("type") || d.HasChange("external_id") || d.HasChange("display_name") {
		org.Name = d.Get("name").(string)
		org.Description = d.Get("description").(string)
		org.Type = d.Get("type").(string)
		org.ExternalID = d.Get("external_id").(string)
		org.DisplayName = d.Get("display_name").(string)
		_, _, err = client.Organizations.UpdateOrganization(*org)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}

func resourceIAMOrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	waitForDelete := d.Get("wait_for_delete").(bool)

	ok, _, err := client.Organizations.DeleteOrganization(*org)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	if waitForDelete {
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"IN_PROGRESS", "QUEUED", "indeterminate"},
			Target:     []string{"SUCCESS"},
			Refresh:    checkOrgDeleteStatus(client, id),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: time.Duration(5) * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("waiting for delete: %w", err))
		}
	}
	return diags
}

func checkOrgDeleteStatus(client *iam.Client, id string) retry.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		orgStatus, resp, err := client.Organizations.DeleteStatus(id)
		if err != nil {
			return resp, "FAILED", err
		}
		if orgStatus != nil {
			return resp, orgStatus.Status, nil
		}
		// We may need to return an error here
		return resp, "IN_PROGRESS", nil
	}
}
