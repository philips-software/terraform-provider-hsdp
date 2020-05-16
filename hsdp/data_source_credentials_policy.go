package hsdp

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	creds "github.com/philips-software/go-hsdp-api/credentials"
)

func dataSourceCredentialsPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCredentialsPolicyRead,
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
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
						},
						"id": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"filter.managing_org", "filter.group_name"},
						},
						"managing_org": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"filter.id"},
						},
						"group_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"filter.id"},
						},
					},
				},
			},
			"policies": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceCredentialsPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	productKey := ""
	managingOrg := ""
	groupName := ""
	id := 0

	if v, ok := d.GetOk("filter"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			groupName = mVi["group_name"].(string)
			managingOrg = mVi["managing_org"].(string)
			productKey = mVi["product_key"].(string)
			id, _ = strconv.Atoi(mVi["id"].(string))
		}
		//managingOrg := d.Get("managing_org").(string)
		//groupName := d.Get("group_name").(string)
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	groupNamePtr := &groupName
	managingOrgPtr := &managingOrg
	idPtr := &id

	if *groupNamePtr == "" {
		groupNamePtr = nil
	}
	if *managingOrgPtr == "" {
		managingOrgPtr = nil
	}
	if id == 0 {
		idPtr = nil
	}
	client, err := config.CredentialsClient()
	if err != nil {
		return err
	}
	if username != "" {
		client, err = config.CredentialsClientWithLogin(username, password)
		if err != nil {
			return err
		}
	}

	creds, _, err := client.Policy.GetPolicy(&creds.GetPolicyOptions{
		ProductKey:  &productKey,
		ManagingOrg: managingOrgPtr,
		GroupName:   groupNamePtr,
		ID:          idPtr,
	})
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(&creds)
	if err != nil {
		return err
	}
	d.SetId("policies")
	d.Set("policies", string(jsonBytes))

	return err
}
