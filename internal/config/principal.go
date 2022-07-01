package config

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Principal represents a HSDP IAM Principal
type Principal struct {
	Username          string
	Password          string
	OAuth2ClientID    string
	OAuth2Password    string
	ServiceID         string
	ServicePrivateKey string
	Region            string
	Environment       string
	Endpoint          string
}

func PrincipalSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func SchemaToPrincipal(d *schema.ResourceData) *Principal {
	principal := Principal{}
	found := false
	if v, ok := d.GetOk("hsdp_principal"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			found = true
			mVi := vi.(map[string]interface{})
			principal.Endpoint = mVi["endpoint"].(string)
			principal.Username = mVi["username"].(string)
			principal.Password = mVi["password"].(string)
			principal.ServiceID = mVi["service_id"].(string)
			principal.ServicePrivateKey = mVi["service_private_key"].(string)
			principal.Environment = mVi["environment"].(string)
			principal.Region = mVi["region"].(string)
			principal.OAuth2ClientID = mVi["oauth2_client_id"].(string)
			principal.OAuth2Password = mVi["oauth2_password"].(string)
		}
	}
	if !found {
		return nil
	}
	return &principal
}
