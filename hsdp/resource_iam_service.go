package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"
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
			"validity": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      12,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 600),
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
				MinItems: 1, // openid
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_scopes": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 1, // openid
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIAMServiceCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	var s iam.Service
	s.Description = d.Get("description").(string)
	s.Name = d.Get("name").(string)
	s.ApplicationID = d.Get("application_id").(string)
	s.Validity = d.Get("validity").(int)
	scopes := expandStringList(d.Get("scopes").(*schema.Set).List())
	defaultScopes := expandStringList(d.Get("default_scopes").(*schema.Set).List())

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
	d.Set("description", s.Description) // RITM0021326

	// Set scopes and default_scopes
	client.Services.AddScopes(*createdService, scopes, defaultScopes)

	return nil
}

func resourceIAMServiceRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	s, resp, err := client.Services.GetServiceByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	// Until RITM0021326 is implemented, this will always clear the field
	// d.Set("description", s.Description)

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
	config := m.(*Config)
	client := config.IAMClient()

	var s iam.Service
	s.ID = d.Id()

	d.Partial(true)
	if d.HasChange("scopes") {
		o, n := d.GetChange("scopes")
		old := expandStringList(o.(*schema.Set).List())
		new := expandStringList(n.(*schema.Set).List())
		toAdd := difference(new, old)
		toRemove := difference(old, new)
		if len(toRemove) > 0 {
			_, _, err := client.Services.RemoveScopes(s, toRemove, []string{})
			if err != nil {
				return err
			}
		}
		if len(toAdd) > 0 {
			client.Services.AddScopes(s, toAdd, []string{})
		}
		d.SetPartial("scopes")
	}
	if d.HasChange("default_scopes") {
		o, n := d.GetChange("default_scopes")
		old := expandStringList(o.(*schema.Set).List())
		new := expandStringList(n.(*schema.Set).List())
		toAdd := difference(new, old)
		toRemove := difference(old, new)
		if len(toRemove) > 0 {
			_, _, err := client.Services.RemoveScopes(s, []string{}, toRemove)
			if err != nil {
				return err
			}
		}
		if len(toAdd) > 0 {
			client.Services.AddScopes(s, []string{}, toAdd)
		}
		d.SetPartial("default_scopes")
	}
	return nil
}

func resourceIAMServiceDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	var s iam.Service
	s.ID = d.Id()
	ok, _, err := client.Services.DeleteService(s)
	if !ok {
		return err
	}
	d.SetId("")
	return nil
}
