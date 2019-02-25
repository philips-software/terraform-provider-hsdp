package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
)

func dataSourceIAMPermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMPermissionsRead,
		Schema: map[string]*schema.Schema{
			"permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}

}

func dataSourceIAMPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*iam.Client)

	resp, _, err := client.Permissions.GetPermissions(nil) // Get all permissions

	if err != nil {
		return err
	}
	permissions := make([]string, 0)

	for _, p := range *resp {
		permissions = append(permissions, p.Name)
	}
	d.SetId("permissions")
	d.Set("permissions", permissions)

	return err
}
