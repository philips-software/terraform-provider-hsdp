package hsdp

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMClient() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMClientCreate,
		Read:   resourceIAMClientRead,
		Update: resourceIAMClientUpdate,
		Delete: resourceIAMClientDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"client_id": &schema.Schema{
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: suppressCaseDiffs,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Id() != "" {
						return true
					}
					return false
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"application_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"global_reference_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"redirection_uris": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"response_types": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"scopes": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_scopes": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"consent_implied": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"access_token_lifetime": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"refresh_token_lifetime": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2592000,
			},
			"id_token_lifetime": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  31536000,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceIAMClientCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var cl iam.ApplicationClient
	cl.Description = d.Get("description").(string)
	cl.Name = d.Get("name").(string)
	cl.ClientID = d.Get("client_id").(string)
	cl.Type = d.Get("type").(string)
	cl.GlobalReferenceID = d.Get("global_reference_id").(string)
	cl.Password = d.Get("password").(string)
	cl.Name = d.Get("name").(string)
	cl.RedirectionURIs = expandStringList(d.Get("redirection_uris").(*schema.Set).List())
	cl.ResponseTypes = expandStringList(d.Get("response_types").(*schema.Set).List())
	cl.ApplicationID = d.Get("application_id").(string)
	cl.Scopes = expandStringList(d.Get("scopes").(*schema.Set).List())
	cl.DefaultScopes = expandStringList(d.Get("default_scopes").(*schema.Set).List())
	cl.IDTokenLifetime = d.Get("id_token_lifetime").(int)
	cl.RefreshTokenLifetime = d.Get("refresh_token_timeline").(int)
	cl.AccessTokenLifetime = d.Get("access_token_lifetime").(int)

	createdClient, _, err := client.Clients.CreateClient(cl)
	if err != nil {
		return err
	}
	d.SetId(createdClient.ID)
	_ = d.Set("password", cl.Password)
	return nil
}

func resourceIAMClientRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	id := d.Id()
	cl, resp, err := client.Clients.GetClientByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("description", cl.Description)
	_ = d.Set("name", cl.Name)
	_ = d.Set("client_id", cl.ClientID)
	_ = d.Set("password", cl.Password)
	_ = d.Set("type", cl.Type)
	_ = d.Set("application_id", cl.ApplicationID)
	_ = d.Set("global_reference_id", cl.GlobalReferenceID)
	_ = d.Set("redirection_uris", cl.RedirectionURIs)
	_ = d.Set("response_types", cl.ResponseTypes)
	_ = d.Set("scopes", cl.Scopes)
	_ = d.Set("default_scopes", cl.DefaultScopes)
	_ = d.Set("access_token_lifetime", cl.AccessTokenLifetime)
	_ = d.Set("refresh_token_lifetime", cl.RefreshTokenLifetime)
	_ = d.Set("id_token_lifetime", cl.IDTokenLifetime)
	_ = d.Set("disabled", cl.Disabled)
	_ = d.Set("consent_implied", cl.ConsentImplied)

	return nil
}

func resourceIAMClientUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var cl iam.ApplicationClient
	cl.ID = d.Id()

	if d.HasChange("scopes") || d.HasChange("default_scopes") {
		newScopes := expandStringList(d.Get("scopes").(*schema.Set).List())
		newDefaultScopes := expandStringList(d.Get("default_scopes").(*schema.Set).List())
		if d.HasChange("scopes") {
			_, ns := d.GetChange("scopes")
			newScopes = expandStringList(ns.(*schema.Set).List())
		}
		if d.HasChange("default_scopes") {
			_, nd := d.GetChange("default_scopes")
			newDefaultScopes = expandStringList(nd.(*schema.Set).List())
		}
		_, _, err := client.Clients.UpdateScopes(cl, newScopes, newDefaultScopes)
		return err
	}
	if d.HasChange("access_token_lifetime") ||
		d.HasChange("refresh_token_lifetime") ||
		d.HasChange("id_token_lifetime") ||
		d.HasChange("consent_implied") ||
		d.HasChange("global_reference_id") ||
		d.HasChange("response_types") ||
		d.HasChange("redirection_uris") {
		cl, _, err := client.Clients.GetClientByID(d.Id())
		if err != nil {
			return err
		}
		cl.RedirectionURIs = expandStringList(d.Get("redirection_uris").(*schema.Set).List())
		cl.ResponseTypes = expandStringList(d.Get("response_types").(*schema.Set).List())
		cl.ConsentImplied = d.Get("consent_implied").(bool)
		cl.AccessTokenLifetime = d.Get("access_token_lifetime").(int)
		cl.RefreshTokenLifetime = d.Get("refresh_token_lifetime").(int)
		cl.IDTokenLifetime = d.Get("id_token_lifetime").(int)
		cl.GlobalReferenceID = d.Get("global_reference_id").(string)
		_, _, err = client.Clients.UpdateClient(*cl)
		if err != nil {
			return err
		}
		return resourceIAMClientRead(d, m)
	}
	return nil
}

func resourceIAMClientDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.IAMClient()
	if err != nil {
		return err
	}

	var cl iam.ApplicationClient
	cl.ID = d.Id()
	ok, _, err := client.Clients.DeleteClient(cl)
	if !ok {
		return err
	}
	d.SetId("")
	return nil
}
