package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIAMOrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMOrgCreate,
		ReadContext:   resourceIAMOrgRead,
		UpdateContext: resourceIAMOrgUpdate,
		DeleteContext: resourceIAMOrgDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"distinct_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_root_org": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"parent_org_id"},
			},
			"parent_org_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"is_root_org"},
			},
		},
	}
}

func resourceIAMOrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isRootOrg := d.Get("is_root_org").(bool)
	externalID := d.Get("external_id").(string)
	orgType := d.Get("type").(string)
	if isRootOrg {
		return diag.FromErr(ErrCannotCreateRootOrg)
	}
	parentOrgID, ok := d.Get("parent_org_id").(string)
	if !ok {
		return diag.FromErr(ErrMissingParentOrgID)
	}
	var newOrg iam.Organization
	newOrg.Name = name
	newOrg.Description = description
	newOrg.Parent.Value = parentOrgID
	newOrg.ExternalID = externalID
	newOrg.Type = orgType
	org, resp, err := client.Organizations.CreateOrganization(newOrg)
	if err != nil {
		return diag.FromErr(err)
	}
	if org == nil {
		return diag.FromErr(fmt.Errorf("failed to create organization: %d", resp.StatusCode))
	}
	d.SetId(org.ID)
	return diags
}

func resourceIAMOrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	org, resp, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	_ = d.Set("org_id", org.ID)
	_ = d.Set("description", org.Description)
	_ = d.Set("name", org.Name)
	_ = d.Set("external_id", org.ExternalID)
	_ = d.Set("parent_org_id", org.Parent.Value)
	_ = d.Set("display_name", org.DisplayName)
	_ = d.Set("active", org.Active)
	_ = d.Set("type", org.Type)
	return diags
}

func resourceIAMOrgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		org.Description = description
		_, _, err = client.Organizations.UpdateOrganization(*org)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}

func resourceIAMOrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ok, _, err := client.Organizations.DeleteOrganization(*org)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(ErrInvalidResponse)
	}
	return diags
}
