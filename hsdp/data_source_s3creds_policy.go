package hsdp

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	creds "github.com/philips-software/go-hsdp-api/s3creds"
)

func dataSourceS3CredsPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceS3CredsPolicyRead,
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
			"product_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"managing_org": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Optional: true,
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

func dataSourceS3CredsPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics
	productKey := ""
	managingOrg := ""
	groupName := ""
	id := 0

	productKey = d.Get("product_key").(string)

	if v, ok := d.GetOk("filter"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			groupName = mVi["group_name"].(string)
			managingOrg = mVi["managing_org"].(string)
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
	client, err := config.S3CredsClient()
	if err != nil {
		return diag.FromErr(err)
	}
	if username != "" {
		client, err = config.S3CredsClientWithLogin(username, password)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	credentials, _, err := client.Policy.GetPolicy(&creds.GetPolicyOptions{
		ProductKey:  &productKey,
		ManagingOrg: managingOrgPtr,
		GroupName:   groupNamePtr,
		ID:          idPtr,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	jsonBytes, err := json.Marshal(&credentials)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("policies")
	_ = d.Set("policies", string(jsonBytes))

	return diags
}
