package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMMFAPolicy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMMFAPolicyCreate,
		ReadContext:   resourceIAMMFAPolicyRead,
		UpdateContext: resourceIAMMFAPolicyUpdate,
		DeleteContext: resourceIAMMFAPolicyDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"user": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"organization"},
				ForceNew:      true,
			},
			"organization": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user"},
				ForceNew:      true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMMFAPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	otpType := d.Get("type").(string)
	active := d.Get("active").(bool)
	user := d.Get("user").(string)
	organization := d.Get("organization").(string)

	if user != "" && organization != "" {
		return diag.FromErr(fmt.Errorf("user and organization are mutually exclusive"))
	}

	var policy iam.MFAPolicy

	policy.Name = name
	policy.Description = description
	policy.SetType(otpType)
	policy.SetActive(active)
	policy.SetResourceOrganization(organization)

	if user != "" {
		policy.SetResourceUser(user)
	}

	newPolicy, resp, err := client.MFAPolicies.CreateMFAPolicy(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	if newPolicy == nil {
		return diag.FromErr(fmt.Errorf("failed to create MFA policy: %d", resp.StatusCode))
	}
	d.SetId(newPolicy.ID)
	_ = d.Set("name", newPolicy.Name)
	_ = d.Set("description", newPolicy.Description)
	_ = d.Set("active", *newPolicy.Active)
	_ = d.Set("version", newPolicy.Meta.Version)
	return diags
}

func resourceIAMMFAPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	policy, resp, err := client.MFAPolicies.GetMFAPolicyByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("description", policy.Description)
	_ = d.Set("active", *policy.Active)
	_ = d.Set("type", policy.Types[0])
	switch policy.Resource.Type {
	case "User":
		_ = d.Set("user", policy.Resource.Value)
	case "Organization":
		_ = d.Set("organization", policy.Resource.Value)
	}
	_ = d.Set("version", policy.Meta.Version)

	return diags
}

func resourceIAMMFAPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	policy, _, err := client.MFAPolicies.GetMFAPolicyByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		policy.Description = description
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		policy.Name = name
	}
	if d.HasChange("active") {
		active := d.Get("active").(bool)
		policy.Active = &active
	}
	updatedPolicy, _, err := client.MFAPolicies.UpdateMFAPolicy(policy)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if updatedPolicy != nil {
		_ = d.Set("version", updatedPolicy.Meta.Version)
	}
	return diags
}

func resourceIAMMFAPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var policy iam.MFAPolicy
	policy.ID = d.Id()

	ok, _, err := client.MFAPolicies.DeleteMFAPolicy(policy)
	if !ok {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
