package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"iam_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["iam_url"],
			},
			"idm_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["idm_url"],
			},
			"oauth2_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["oauth2_client_id"],
			},
			"oauth2_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: descriptions["oauth2_password"],
			},
			"org_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["org_id"],
			},
			"org_admin_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["org_admin_username"],
			},
			"org_admin_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: descriptions["org_admin_password"],
			},
			"shared_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
				Description: descriptions["shared_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: descriptions["secret_key"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hsdp_iam_org":         resourceIAMOrg(),
			"hsdp_iam_group":       resourceIAMGroup(),
			"hsdp_iam_permission":  resourceIAMPermission(),
			"hsdp_iam_role":        resourceIAMRole(),
			"hsdp_iam_proposition": resourceIAMProposition(),
			"hsdp_iam_application": resourceIAMApplication(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hsdp_iam_introspect": dataSourceIAMIntrospect(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"iam_url":            "The HSDP IAM instance URL",
		"idm_url":            "The HSDP IDM instance URL",
		"oauth2_client_id":   "The OAuth2 client id",
		"oauth2_password":    "The OAuth2 password",
		"org_id":             "The (top level) Organization ID - UUID",
		"org_admin_username": "The username of the Organization Admin",
		"org_admin_password": "The password of the Organization Admin",
		"shared_key":         "The shared key",
		"secret_key":         "The secret key",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{}

	config.IAMURL = d.Get("iam_url").(string)
	config.IDMURL = d.Get("idm_url").(string)
	config.OAuth2ClientID = d.Get("oauth2_client_id").(string)
	config.OAuth2Secret = d.Get("oauth2_password").(string)
	config.RootOrgID = d.Get("org_id").(string)
	config.OrgAdminUsername = d.Get("org_admin_username").(string)
	config.OrgAdminPassword = d.Get("org_admin_password").(string)
	config.SharedKey = d.Get("shared_key").(string)
	config.SecretKey = d.Get("secret_key").(string)

	return config.Client()
}
