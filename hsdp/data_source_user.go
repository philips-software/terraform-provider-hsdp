package hsdp

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.IAMClient()

	username := d.Get("username").(string)

	uuid, _, err := client.Users.GetUserIDByLoginID(username)

	if err != nil {
		return err
	}
	d.SetId(uuid)
	d.Set("uuid", uuid)

	return nil
}
