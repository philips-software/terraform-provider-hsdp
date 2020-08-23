package hsdp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"testing"
)

var datasourcConfigTest = `
provider "hsdp" {
    region = "us-east"
    environment = "client-test"
}

data "hsdp_config" "test" {
    region = "us-east"
    environment = "client-test"
	service = "iam"
}
`

func TestAccDataSourceHSDPConfig(t *testing.T) {
	resourceName := "data.hsdp_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourcConfigTest,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", "https://iam-client-test.us-east.philips-healthsuite.com"),
				),
			},
		},
	})
}
