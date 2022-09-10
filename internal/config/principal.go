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
	UAAUsername       string
	UAAPassword       string
}

func PrincipalSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"username": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"password": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"oauth2_client_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"oauth2_password": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
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
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"endpoint": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"uaa_username": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"uaa_password": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
		},
	}
}

func (p *Principal) HasAuth() bool {
	// Service identity
	if p.ServiceID != "" && p.ServicePrivateKey != "" {
		return true
	}
	// IAM identity
	if p.Username != "" && p.Password != "" {
		return true
	}
	if p.UAAUsername != "" && p.UAAPassword != "" {
		return true
	}
	// No credentials
	return false
}

func SchemaToPrincipal(d *schema.ResourceData, m interface{}) *Principal {
	config := m.(*Config)

	principal := Principal{}
	if v, ok := d.GetOk("principal"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		mVi := v.([]interface{})[0].(map[string]interface{})
		principal.Endpoint = mVi["endpoint"].(string)
		principal.Username = mVi["username"].(string)
		principal.Password = mVi["password"].(string)
		principal.ServiceID = mVi["service_id"].(string)
		principal.ServicePrivateKey = mVi["service_private_key"].(string)
		principal.Environment = mVi["environment"].(string)
		principal.Region = mVi["region"].(string)
		principal.OAuth2ClientID = mVi["oauth2_client_id"].(string)
		principal.OAuth2Password = mVi["oauth2_password"].(string)
		principal.UAAUsername = mVi["uaa_username"].(string)
		principal.UAAPassword = mVi["uaa_password"].(string)
	}
	// Set defaults
	if principal.Environment == "" {
		principal.Environment = config.Environment
	}
	if principal.Region == "" {
		principal.Region = config.Region
	}
	return &principal
}
