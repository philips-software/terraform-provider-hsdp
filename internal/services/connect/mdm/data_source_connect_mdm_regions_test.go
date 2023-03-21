package mdm_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDataSourceMDMRegions_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.hsdp_connect_mdm_regions.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccDataSourceMDMRegions(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "ids.0", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("empty id")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith(resourceName, "names.0", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("empty name")
						}
						return nil
					}),
				),
			},
		},
	})
}

func testAccDataSourceMDMRegions() string {
	return `data "hsdp_connect_mdm_regions" "test" {}`
}
