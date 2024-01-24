package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

var descriptions = map[string]string{
	"user": "Users are created under an organization and are assigned to groups.",
}

func ResourceIAMUser() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["user"],
		Importer: &schema.ResourceImporter{
			StateContext: importUserContext,
		},

		CreateContext: resourceIAMUserCreate,
		ReadContext:   resourceIAMUserRead,
		UpdateContext: resourceIAMUserUpdate,
		DeleteContext: resourceIAMUserDelete,

		SchemaVersion: 2,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use login field instead",
			},
			"login": {
				Type:             schema.TypeString,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
				Required:         true,
			},
			"email": {
				Type:             schema.TypeString,
				DiffSuppressFunc: tools.SuppressCaseDiffs,
				Required:         true,
				Description:      "The email address of the user.",
			},
			"password": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				Description: "When specified this will skip the email activation flow and immediately activate the IAM account. Very Important: you are responsible for sharing this password with the new IAM user through some channel of communication. No email will be triggered by the system. If unsure, do not set a password so the normal email activation flow is followed. Finally, any password value changes after user creation will have no effect on the users' actual password.",
			},
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The first name of the user.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The last name of the user.",
			},
			"mobile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional mobile phone number of the user.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managing organization of the user.",
			},
			"preferred_language": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.SuppressEmptyPreferredLanguage,
				Description:      "Language preference for all communications. Value can be a two letter language code as defined by ISO 639-1 (en, de) or it can be a combination of language code and country code (en-gb, en-us). The country code is as per ISO 3166 two letter code (alpha-2).",
			},
			"preferred_communication_channel": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.SuppressDefaultCommunicationChannel,
				Description:      "Preferred communication channel. Email and SMS are supported channels. Email is the default channel if e-mail address is provided. Values supported: [ email | sms ].",
			},
		},
	}
}

func resourceIAMUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	last := d.Get("last_name").(string)
	first := d.Get("first_name").(string)
	email := d.Get("username").(string) // Deprecated
	mobile := d.Get("mobile").(string)
	login := d.Get("login").(string)
	password := d.Get("password").(string)
	reversedPassword := tools.ReverseString(password)
	if login == "" {
		login = email
	}
	email = d.Get("email").(string)
	organization := d.Get("organization_id").(string)
	preferredLanguage := d.Get("preferred_language").(string)
	preferredCommunicationChannel := d.Get("preferred_communication_channel").(string)

	// First check if this user already exists
	foundUser, _, err := client.Users.GetUserByID(login)
	if err == nil && (foundUser != nil && foundUser.ID != "") {
		if foundUser.ManagingOrganization != organization {
			return diag.FromErr(fmt.Errorf("user '%s' already exists but is managed by a different IAM organization", login))
		}
		d.SetId(foundUser.ID)
		return resourceIAMUserRead(ctx, d, m)
	}
	initialPassword := password != "" && client.HasSigningKeys()

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
		ManagingOrganization:          organization,
		PreferredLanguage:             preferredLanguage,
		PreferredCommunicationChannel: preferredCommunicationChannel,
		IsAgeValidated:                "true",
	}
	if initialPassword { // We first use the reverse
		person.Password = reversedPassword
	}
	if mobile != "" {
		person.Telecom = append(person.Telecom,
			iam.TelecomEntry{
				System: "mobile",
				Value:  mobile,
			})
	}
	user, resp, err := client.Users.CreateUser(person)
	if err != nil {
		return diag.FromErr(err)
	}
	if user == nil {
		return diag.FromErr(fmt.Errorf("error creating user '%s': %v %w", login, resp, err))
	}
	if initialPassword { // Set the final password
		success, _, err := client.Users.ChangePassword(login, reversedPassword, password)
		if !success {
			person.ID = user.ID
			_, _, _ = client.Users.DeleteUser(person)
			return diag.FromErr(fmt.Errorf("error setting password for '%s': %w", login, err))
		}
	}
	d.SetId(user.ID)
	return resourceIAMUserRead(ctx, d, m)
}

func resourceIAMUserRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()

	user, _, err := client.Users.GetUserByID(id)
	if err != nil {
		if errors.Is(err, iam.ErrEmptyResults) {
			// Means the user was cleared, probably due to not activating their account
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
	_ = d.Set("preferred_communication_channel", user.PreferredCommunicationChannel)
	_ = d.Set("preferred_language", user.PreferredLanguage)
	return diags
}

func resourceIAMUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.IAMClient()
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
	if d.HasChange("last_name") || d.HasChange("first_name") || d.HasChange("email") ||
		d.HasChange("mobile") || d.HasChange("preferred_language") || d.HasChange("preferred_communication_channel") {
		profile, _, err := client.Users.LegacyGetUserByUUID(d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("resourceIAMUserUpdate LegacyGetUserByUUID: %w", err))
		}
		profile.FamilyName = d.Get("last_name").(string)
		profile.GivenName = d.Get("first_name").(string)
		profile.PreferredLanguage = d.Get("preferred_language").(string)
		profile.PreferredCommunicationChannel = d.Get("preferred_communication_channel").(string)
		profile.Contact.EmailAddress = d.Get("email").(string)
		if profile.MiddleName == "" {
			profile.MiddleName = " "
		}
		profile.Contact.MobilePhone = d.Get("mobile").(string)
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
	readDiags := resourceIAMUserRead(ctx, d, m)
	return append(diags, readDiags...)
}

func resourceIAMUserDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.IAMClient()
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
	_, resp, err := client.Users.DeleteUser(person)
	if err != nil {
		return diag.FromErr(fmt.Errorf("DeleteUser '%s' error: %w", person.ID, err))
	}
	if resp != nil && resp.StatusCode() == http.StatusConflict {
		return diag.FromErr(fmt.Errorf("DeleteUser return HTTP 409 Conflict: %w", err))
	}
	if resp != nil && resp.StatusCode() != http.StatusNoContent {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "DeleteUser returned unexpected result",
			Detail:   fmt.Sprintf("DeleteUser returned status '%d', which is unexpected: %v", resp.StatusCode(), err),
		})
	}
	d.SetId("")
	return diags
}
