package role_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMRole_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_role.test"
	parentOrgID := acc.AccIAMOrgGUID()
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMRole(parentOrgID, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "managing_organization", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMRole(parentOrgID, name string) string {
	roleName := fmt.Sprintf("TESTROLE-%s", strings.ToUpper(name))
	return fmt.Sprintf(`
resource "hsdp_iam_role" "test" {
  name        = "%s"
  description = "Acceptance Test Role %s"

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
  managing_organization = "%s"
}`, roleName, roleName, parentOrgID)
}
