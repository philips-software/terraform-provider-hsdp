package role_sharing_policy_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDatasourceIAMRoleSharingPolicies_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.hsdp_iam_role_sharing_policies.test"
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
				Config:       testAccDatasourceIAMRoleSharingPolicies(parentOrgID, randomName, "AllowChildren"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sharing_policies.0", "AllowChildren"),
					resource.TestCheckResourceAttr(resourceName, "purposes.0", "ACCTEST"),
				),
			},
		},
	})
}

func testAccDatasourceIAMRoleSharingPolicies(parentOrgID, name, sharingPolicy string) string {
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

data "hsdp_iam_role_sharing_policies" "test" {
  role_id                = hsdp_iam_role.test.id

  depends_on = [hsdp_iam_role_sharing_policy.test]
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
