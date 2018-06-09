package hsdp

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hsdp/go-hsdp-api/iam"
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
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
	client := m.(*iam.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isRootOrg := d.Get("is_root_org").(bool)
	if isRootOrg {
		return errors.New("cannot create root orgs")
	}
	parentOrgID, ok := d.Get("parent_org_id").(string)
	if !ok {
		return errors.New("non root orgs must specify a `parent_org_id`")
	}
	org, resp, err := client.Organizations.CreateOrganization(parentOrgID, name, description)
	if err != nil {
		return err
	}
	if org == nil {
		return fmt.Errorf("failed to create organization: %d", resp.StatusCode)
	}
	d.SetId(org.OrganizationID)
	return nil
}

func resourceIAMOrgRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return err
	}
	d.Set("org_id", org.OrganizationID)
	d.Set("description", org.Description)
	d.Set("name", org.Name)
	return nil
}

func resourceIAMOrgUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return err
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		org.Description = description
		client.Organizations.UpdateOrganization(*org)
	}
	return nil
}

func resourceIAMOrgDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
