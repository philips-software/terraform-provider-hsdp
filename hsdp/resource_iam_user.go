package hsdp

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMUser() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Create: resourceIAMUserCreate,
		Read:   resourceIAMUserRead,
		Update: resourceIAMUserUpdate,
		Delete: resourceIAMUserDelete,

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
				Required: true,
			},
		},
	}
}

func resourceIAMUserCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	last := d.Get("last_name").(string)
	first := d.Get("first_name").(string)
	email := d.Get("username").(string) // Deprecated
	mobile := d.Get("mobile").(string)
	login := d.Get("login").(string)
	if login == "" {
		login = email
	}
	email = d.Get("email").(string)
	organization := d.Get("organization_id").(string)

	// First check if this user already exists
	uuid, _, err := client.Users.GetUserIDByLoginID(email)
	if err == nil && uuid != "" {
		user, _, _ := client.Users.GetUserByID(uuid)
		if user != nil {
			if user.AccountStatus.Disabled {
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
	user, _, err := client.Users.CreateUser(person)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("Error creating user")
	}
	d.SetId(user.ID)
	return nil
}

func resourceIAMUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()

	user, _, err := client.Users.GetUserByID(id)
	if err != nil {
		if _, ok := err.(*iam.UserError); ok {
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("login", user.LoginID)
	_ = d.Set("last_name", user.Name.Family)
	_ = d.Set("first_name", user.Name.Given)
	_ = d.Set("email", user.EmailAddress)
	_ = d.Set("login", user.LoginID)
	_ = d.Set("organization_id", user.ManagingOrganization)
	return nil
}

func resourceIAMUserUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var p iam.Person
	p.ID = d.Id()

	if d.HasChange("login") {
		newLogin := d.Get("login").(string)
		_, _, err := client.Users.ChangeLoginID(p, newLogin)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIAMUserDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()

	user, _, err := client.Users.GetUserByID(id)
	if err != nil {
		if _, ok := err.(*iam.UserError); ok {
			d.SetId("")
			return nil
		}
		return err
	}
	if user == nil {
		return nil
	}
	var person iam.Person
	person.ID = user.ID
	ok, _, _ := client.Users.DeleteUser(person)
	if ok {
		d.SetId("")
	}
	return nil
}
