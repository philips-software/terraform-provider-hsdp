package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDataSourceIAMIntrospect_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.hsdp_iam_introspect.test"
	org := acc.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccDataSourceIAMIntrospect(org),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "effective_permissions.0", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("empty set")
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "4"),
				),
			},
		},
	})
}

func testAccDataSourceIAMIntrospect(org string) string {
	return fmt.Sprintf(`
data "hsdp_iam_introspect" "test" {
    organization_context = "%s"
}`, org)
}
