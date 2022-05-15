package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acctest"
)

func TestAccResourceIAMRole_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_role.test"
	parentOrgID := acctest.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMRole(parentOrgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "managing_organization", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMRole(parentOrgID string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_role" "test" {
  name        = "TESTROLE"
  description = "Acceptance Test Role"

  permissions = [
    "DATAITEM.CREATEONBEHALF",
    "DATAITEM.READ",
    "DATAITEM.DELETEONBEHALF",
    "DATAITEM.DELETE",
    "CONTRACT.CREATE",
    "DATAITEM.READONBEHALF",
    "CONTRACT.READ",
    "DATAITEM.CREATE",
  ]
  managing_organization = %[1]q
}`, parentOrgID)
}
