package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
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
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Id() != "" {
						return true
					}
					return false
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": &schema.Schema{
				Type:     schema.TypeString,
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
		},
	}
}

func resourceIAMClientCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
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

	createdClient, _, err := client.Clients.CreateClient(cl)
	if err != nil {
		return err
	}
	d.SetId(createdClient.ID)
	return nil
}

func resourceIAMClientRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)

	id := d.Id()
	cl, _, err := client.Clients.GetClientByID(id)
	if err != nil {
		return err
	}
	d.Set("description", cl.Description)
	d.Set("name", cl.Name)
	d.Set("client_id", cl.ClientID)
	d.Set("password", cl.Password)
	d.Set("type", cl.Type)
	d.Set("application_id", cl.ApplicationID)
	d.Set("global_reference_id", cl.GlobalReferenceID)
	d.Set("redirection_uris", cl.RedirectionURIs)
	d.Set("response_types", cl.ResponseTypes)
	return nil
}

func resourceIAMClientUpdate(d *schema.ResourceData, m interface{}) error {
	//client := m.(*iam.Client)
	var cl iam.ApplicationClient
	cl.ID = d.Id()
	// TODO
	return nil
}

func resourceIAMClientDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*iam.Client)
	var cl iam.ApplicationClient
	cl.ID = d.Id()
	ok, _, err := client.Clients.DeleteClient(cl)
	if !ok {
		return err
	}
	d.SetId("")
	return nil
}
