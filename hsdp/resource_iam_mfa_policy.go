package hsdp

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMMFAPolicy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMMFAPolicyCreate,
		Read:   resourceIAMMFAPolicyRead,
		Update: resourceIAMMFAPolicyUpdate,
		Delete: resourceIAMMFAPolicyDelete,

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

func resourceIAMMFAPolicyCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	otpType := d.Get("type").(string)
	active := d.Get("active").(bool)
	user := d.Get("user").(string)
	organization := d.Get("organization").(string)

	if user != "" && organization != "" {
		return fmt.Errorf("user and organization are mutually exclusive")
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
		return err
	}
	if newPolicy == nil {
		return fmt.Errorf("failed to create MFA policy: %d", resp.StatusCode)
	}
	d.SetId(newPolicy.ID)
	d.Set("name", newPolicy.Name)
	d.Set("description", newPolicy.Description)
	d.Set("active", *newPolicy.Active)
	d.Set("version", newPolicy.Meta.Version)
	return nil
}

func resourceIAMMFAPolicyRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	policy, resp, err := client.MFAPolicies.GetMFAPolicyByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("active", *policy.Active)
	d.Set("type", policy.Types[0])
	switch policy.Resource.Type {
	case "User":
		d.Set("user", policy.Resource.Value)
	case "Organization":
		d.Set("organization", policy.Resource.Value)
	}
	d.Set("version", policy.Meta.Version)

	return nil
}

func resourceIAMMFAPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	policy, _, err := client.MFAPolicies.GetMFAPolicyByID(id)
	if err != nil {
		return err
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
	if updatedPolicy != nil {
		d.Set("version", updatedPolicy.Meta.Version)
	}
	return err
}

func resourceIAMMFAPolicyDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	var policy iam.MFAPolicy
	policy.ID = d.Id()

	ok, _, err := client.MFAPolicies.DeleteMFAPolicy(policy)
	if !ok {
		return err
	}
	d.SetId("")
	return nil
}
