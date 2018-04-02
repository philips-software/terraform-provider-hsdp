package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	iamclient "github.com/loafoe/go-hsdpiam"
)

func resourceIAMGroup() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMGroupCreate,
		Read:   resourceIAMGroupRead,
		Update: resourceIAMGroupUpdate,
		Delete: resourceIAMGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"managing_organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIAMGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iamclient.Client)
	var group iamclient.Group
	group.Description = d.Get("description").(string)
	group.Name = d.Get("name").(string)
	group.ManagingOrganization = d.Get("managing_organization").(string)

	createdGroup, _, err := client.Groups.CreateGroup(group)
	if err != nil {
		return err
	}
	d.SetId(createdGroup.ID)
	d.Set("name", createdGroup.Name)
	d.Set("description", createdGroup.Description)
	d.Set("managing_organization", createdGroup.ManagingOrganization)
	return nil
}

func resourceIAMGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iamclient.Client)

	id := d.Id()
	group, _, err := client.Groups.GetGroup(&iamclient.GetGroupOptions{ID: &id})
	if err != nil {
		return err
	}
	d.Set("managing_organization", group.ManagingOrganization)
	d.Set("description", group.Description)
	d.Set("name", group.Name)
	return nil
}

func resourceIAMGroupUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("description") {
		return nil
	}
	client := m.(*iamclient.Client)
	var group iamclient.Group
	group.ID = d.Id()
	group.Description = d.Get("description").(string)
	_, _, err := client.Groups.UpdateGroup(group)
	if err != nil {
		return err
	}
	return nil
}

func resourceIAMGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*iamclient.Client)
	var group iamclient.Group
	group.ID = d.Id()
	_, _, err := client.Groups.DeleteGroup(group)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
