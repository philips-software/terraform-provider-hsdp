package hsdp

import (
	"fmt"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIAMOrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMOrgCreate,
		Read:   resourceIAMOrgRead,
		Update: resourceIAMOrgUpdate,
		Delete: resourceIAMOrgDelete,

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

func resourceIAMOrgCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isRootOrg := d.Get("is_root_org").(bool)
	externalID := d.Get("external_id").(string)
	orgType := d.Get("type").(string)
	if isRootOrg {
		return ErrCannotCreateRootOrg
	}
	parentOrgID, ok := d.Get("parent_org_id").(string)
	if !ok {
		return ErrMissingParentOrgID
	}
	var newOrg iam.Organization
	newOrg.Name = name
	newOrg.Description = description
	newOrg.Parent.Value = parentOrgID
	newOrg.ExternalID = externalID
	newOrg.Type = orgType
	org, resp, err := client.Organizations.CreateOrganization(newOrg)
	if err != nil {
		return err
	}
	if org == nil {
		return fmt.Errorf("failed to create organization: %d", resp.StatusCode)
	}
	d.SetId(org.ID)
	return nil
}

func resourceIAMOrgRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	org, resp, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("org_id", org.ID)
	_ = d.Set("description", org.Description)
	_ = d.Set("name", org.Name)
	_ = d.Set("external_id", org.ExternalID)
	_ = d.Set("parent_org_id", org.Parent.Value)
	_ = d.Set("display_name", org.DisplayName)
	_ = d.Set("active", org.Active)
	_ = d.Set("type", org.Type)
	return nil
}

func resourceIAMOrgUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return err
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		org.Description = description
		_, _, err = client.Organizations.UpdateOrganization(*org)
	}
	return err
}

func resourceIAMOrgDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return err
	}

	ok, _, err := client.Organizations.DeleteOrganization(*org)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidResponse
	}
	return nil
}
