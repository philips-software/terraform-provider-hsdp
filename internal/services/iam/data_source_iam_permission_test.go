package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDataSourceIAMPermission_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.hsdp_iam_permission.test"
	name := "ALL.WRITE"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccDataSourceIAMPermission(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "category", "NATIVECDRTEST"),
				),
			},
		},
	})
}

func testAccDataSourceIAMPermission(name string) string {
	return fmt.Sprintf(`
data "hsdp_iam_permission" "test" {
    name = "%s"
}`, name)
}
