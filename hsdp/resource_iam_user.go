package hsdp

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hsdp/go-hsdp-api/iam"
)

func resourceIAMUser() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMUserCreate,
		Read:   resourceIAMUserRead,
		Update: resourceIAMUserUpdate,
		Delete: resourceIAMUserDelete,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mobile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIAMUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
	last := d.Get("last_name").(string)
	first := d.Get("first_name").(string)
	email := d.Get("username").(string)
	mobile := d.Get("mobile").(string)

	ok, _, err := client.Users.CreateUser(first, last, email, mobile, "")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Error creating user")
	}
	// Fetch UUID
	uuid, _, err := client.Users.GetUserIDByLoginID(email)
	if err != nil {
		return fmt.Errorf("Cannot find newly minted user")
	}
	d.SetId(uuid)
	return nil
}

func resourceIAMUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)

	id := d.Id()

	user, _, err := client.Users.GetUserByID(id)
	if err != nil {
		return err
	}
	d.Set("last_name", user.Name.Family)
	d.Set("first_name", user.Name.Given)
	for _, t := range user.Telecom {
		if t.System == "email" {
			d.Set("username", t.Value)
			continue
		}
		if t.System == "mobile" {
			d.Set("mobile", t.Value)
			continue
		}
	}
	return nil
}

func resourceIAMUserUpdate(d *schema.ResourceData, m interface{}) error {
	// Not supported by API at this time
	return nil
}

func resourceIAMUserDelete(d *schema.ResourceData, m interface{}) error {
	// Not supported by API at this time
	return nil
}
