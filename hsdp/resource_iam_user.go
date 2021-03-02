package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMUser() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMUserCreate,
		ReadContext:   resourceIAMUserRead,
		UpdateContext: resourceIAMUserUpdate,
		DeleteContext: resourceIAMUserDelete,

		SchemaVersion: 1,
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
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
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

func resourceIAMUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	last := d.Get("last_name").(string)
	first := d.Get("first_name").(string)
	email := d.Get("username").(string) // Deprecated
	mobile := d.Get("mobile").(string)
	login := d.Get("login").(string)
	password := d.Get("password").(string)
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
				return diag.FromErr(err)
			}
			diags = resourceIAMUserRead(ctx, d, m)
			if len(diags) == 0 {
				d.SetId(user.ID)
			}
			return diags
		}
	}
	person := iam.Person{
		ResourceType: "Person",
		Name: iam.Name{
			Family: last,
			Given:  first,
		},
		LoginID:  login,
		Password: password,
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
		return diag.FromErr(err)
	}
	if user == nil {
		return diag.FromErr(fmt.Errorf("Error creating user: %w", err))
	}
	d.SetId(user.ID)
	return diags
}

func resourceIAMUserRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()

	user, _, err := client.Users.GetUserByID(id)
	if err != nil {
		// Means the user was cleared, probably due to not activating their account
		if _, ok := err.(*iam.UserError); ok {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("login", user.LoginID)
	_ = d.Set("last_name", user.Name.Family)
	_ = d.Set("first_name", user.Name.Given)
	_ = d.Set("email", user.EmailAddress)
	_ = d.Set("login", user.LoginID)
	_ = d.Set("organization_id", user.ManagingOrganization)
	return diags
}

func resourceIAMUserUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var p iam.Person
	p.ID = d.Id()

	if d.HasChange("login") {
		newLogin := d.Get("login").(string)
		_, _, err := client.Users.ChangeLoginID(p, newLogin)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("last_name") || d.HasChange("first_name") || d.HasChange("email") {
		profile, _, err := client.Users.LegacyGetUserByUUID(d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("resourceIAMUserUpdate LegacyGetUserByUUID: %w", err))
		}
		profile.FamilyName = d.Get("last_name").(string)
		profile.GivenName = d.Get("first_name").(string)
		profile.Contact.EmailAddress = d.Get("email").(string)
		if profile.MiddleName == "" {
			profile.MiddleName = " "
		}
		profile.ID = d.Id()
		_, _, err = client.Users.LegacyUpdateUser(*profile)
		if err != nil {
			return diag.FromErr(fmt.Errorf("resourceIAMUserUpdate LegacyUpdateUser: %w", err))
		}
	}
	if d.HasChange("password") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "password change not propagated",
			Detail:   "changing the password after a user is created has no effect",
		})
	}
	return diags
}

func resourceIAMUserDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()

	user, _, err := client.Users.GetUserByID(id)
	if err != nil {
		if _, ok := err.(*iam.UserError); ok {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	if user == nil {
		return diags
	}
	var person iam.Person
	person.ID = user.ID
	ok, _, _ := client.Users.DeleteUser(person)
	if ok {
		d.SetId("")
	}
	return diags
}
