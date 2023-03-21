package configuration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDataSourceConfig_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.hsdp_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", "https://iam-client-test.us-east.philips-healthsuite.com"),
				),
			},
		},
	})
}

func testAccDataSourceConfig() string {
	return `
data "hsdp_config" "test" {
  service     = "iam"
  region      = "us-east"
  environment = "client-test"
}`
}
