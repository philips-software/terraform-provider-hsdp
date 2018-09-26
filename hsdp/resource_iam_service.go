package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMService() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMServiceCreate,
		Read:   resourceIAMServiceRead,
		Update: resourceIAMServiceUpdate,
		Delete: resourceIAMServiceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"service_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"scopes": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_scopes": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIAMServiceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
	var s iam.Service
	s.Description = d.Get("description").(string)
	s.Name = d.Get("name").(string)
	s.ApplicationID = d.Get("application_id").(string)

	createdService, _, err := client.Services.CreateService(s)
	if err != nil {
		return err
	}
	d.SetId(createdService.ID)
	d.Set("expires_on", createdService.ExpiresOn)
	d.Set("scopes", createdService.Scopes)
	d.Set("default_scopes", createdService.DefaultScopes)
	d.Set("private_key", createdService.PrivateKey)
	d.Set("service_id", createdService.ServiceID)
	d.Set("organization_id", createdService.OrganizationID)
	return nil
}

func resourceIAMServiceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)

	id := d.Id()
	s, _, err := client.Services.GetServiceByID(id)
	if err != nil {
		return err
	}
	d.Set("description", s.Description)
	d.Set("name", s.Name)
	d.Set("application_id", s.ApplicationID)
	d.Set("expires_on", s.ExpiresOn)
	d.Set("organization_id", s.OrganizationID)
	d.Set("service_id", s.ServiceID)
	d.Set("scopes", s.Scopes)
	d.Set("expires_on", s.ExpiresOn)
	d.Set("default_scopes", s.DefaultScopes)
	// The private key is only returned on create
	return nil
}

func resourceIAMServiceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
	var s iam.Service
	s.ID = d.Id()

	d.Partial(true)
	if d.HasChange("scopes") || d.HasChange("default_scopes") {
		newScopes := expandStringList(d.Get("scopes").(*schema.Set).List())
		newDefaultScopes := expandStringList(d.Get("default_scopes").(*schema.Set).List())
		if d.HasChange("scopes") {
			_, ns := d.GetChange("scopes")
			newScopes = expandStringList(ns.(*schema.Set).List())
			d.SetPartial("scopes")
		}
		if d.HasChange("default_scopes") {
			_, nd := d.GetChange("default_scopes")
			newDefaultScopes = expandStringList(nd.(*schema.Set).List())
			d.SetPartial("default_scopes")
		}
		_, _, err := client.Services.UpdateScopes(s, newScopes, newDefaultScopes)
		return err
	}
	return nil
}

func resourceIAMServiceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
	var s iam.Service
	s.ID = d.Id()
	ok, _, err := client.Services.DeleteService(s)
	if !ok {
		return err
	}
	d.SetId("")
	return nil
}
