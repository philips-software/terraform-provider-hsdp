package client

import (
	"context"
	"net/http"

	"github.com/philips-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceIAMClient() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMClientRead,
		Description: descriptions["client"],
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_reference_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"redirection_uris": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"response_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"consent_implied": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"access_token_lifetime": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"refresh_token_lifetime": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id_token_lifetime": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceIAMClientRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	token, err := client.Token()
	if err != nil || token == "" {
		return diag.FromErr(config.ErrMissingIAMCredentials)
	}

	name := d.Get("name").(string)
	applicationId := d.Get("application_id").(string)

	clients, resp, err := client.Clients.GetClients(&iam.GetClientsOptions{
		Name:          &name,
		ApplicationID: &applicationId,
	})
	if err != nil || clients == nil || len(*clients) == 0 {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	cl := (*clients)[0]

	_ = d.Set("description", cl.Description)
	_ = d.Set("name", cl.Name)
	_ = d.Set("client_id", cl.ClientID)
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
	d.SetId(cl.ID)

	return diags
}
