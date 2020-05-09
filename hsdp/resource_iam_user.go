package hsdp

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
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

		DeprecationMessage: "Please use the HSDP IAM self service portal for user management",

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use login field instead",
			},
			"login": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
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
				Optional: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIAMUserCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	last := d.Get("last_name").(string)
	first := d.Get("first_name").(string)
	email := d.Get("username").(string) // Deprecated
	mobile := d.Get("mobile").(string)
	login := d.Get("login").(string)
	email = d.Get("email").(string)
	organization := d.Get("organization_id").(string)

	// First check if this user already exists
	uuid, _, err := client.Users.GetUserIDByLoginID(email)
	if err == nil && uuid != "" {
		user, _, _ := client.Users.GetUserByID(uuid)
		if user != nil {
			if user.Disabled {
				// Retrigger activation email
				_, _, err = client.Users.ResendActivation(email)
				return err
			}
			err = resourceIAMUserRead(d, m)
			if err == nil {
				d.SetId(user.ID)
			}
			return nil
		}
	}
	person := iam.Person{
		ResourceType: "Person",
		Name: iam.Name{
			Family: last,
			Given:  first,
		},
		LoginID: login,
		Telecom: []iam.TelecomEntry{
			{
				System: "email",
				Value:  email,
			},
		},
		ManagingOrganization: organization,
		IsAgeValidated:       "true",
	}
	if mobile != "" {
		person.Telecom = append(person.Telecom,
			iam.TelecomEntry{
				System: "mobile",
				Value:  mobile,
			})
	}
	ok, _, err := client.Users.CreateUser(person)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Error creating user")
	}
	// Fetch UUID
	uuid, _, err = client.Users.GetUserIDByLoginID(email)
	if err != nil {
		return fmt.Errorf("Cannot find newly minted user")
	}
	d.SetId(uuid)
	return nil
}

func resourceIAMUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()

	user, _, err := client.Users.GetUserByID(id)
	if err != nil {
		if _, ok := err.(*iam.UserError); ok {
			d.SetId("")
			return nil
		}
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
