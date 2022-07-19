package role_sharing_policy_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMRoleSharingPolicy_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_role_sharing_policy.test"
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
				Config:       testAccResourceIAMRoleSharingPolicy(parentOrgID, randomName, "AllowChildren"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sharing_policy", "AllowChildren"),
				),
			},
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMRoleSharingPolicy(parentOrgID, randomName, "Denied"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sharing_policy", "Denied"),
				),
			},
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMRoleSharingPolicy(parentOrgID, randomName, "Restricted"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sharing_policy", "Restricted"),
				),
			},
		},
	})
}

func testAccResourceIAMRoleSharingPolicy(parentOrgID, name, sharingPolicy string) string {
	roleName := fmt.Sprintf("TESTROLE-%s", strings.ToUpper(name))
	return fmt.Sprintf(`
resource "hsdp_iam_org" "another" {
  name = "ACC-%s"
  description = "IAM Role Sharing Policy Test %s"

  parent_org_id = "%s"
  wait_for_delete = false
}

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
}

resource "hsdp_iam_role_sharing_policy" "test" {
  role_id                = hsdp_iam_role.test.id
  purpose                = "ACCTEST"
  target_organization_id = hsdp_iam_org.another.id
  sharing_policy         = "%s"
}
`,
		// ORG
		name,
		name,
		parentOrgID,

		// ROLE
		roleName,
		name,
		parentOrgID,

		// ROLE SHARING POLICY
		sharingPolicy,
	)
}
