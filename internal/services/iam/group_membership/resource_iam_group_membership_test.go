package group_membership_test

import (
	"fmt"
	"testing"

	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMGroupMembership_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_user.test"
	parentOrgID := acc.AccIAMOrgGUID()
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randomPassword, _ := tools.RandomPassword()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMUser(parentOrgID, randomName, randomPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "organization_id", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMUser(parentOrgID, name, password string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_user" "test" {
  login           = "%s"
  email           = "acceptance+%s@terrakube.com"
  first_name      = "ACC"
  last_name       = "Developer"
  password        = "%s"
  organization_id = "%s"
}

resource "hsdp_iam_group" "test" {
  name = "test-%s"
  managing_organization = "%s"
  description = "Acceptance Test for User"
  users = []
  roles = []

  drift_detection = false
}

resource "hsdp_iam_group_membership" "test" {
  iam_group_id = hsdp_iam_group.test.id
  users = [hsdp_iam_user.test.id]
}
`,
		// IAM_USER
		name, name, password, parentOrgID,
		// IAM_GROUP
		name,
		parentOrgID,
		// IAM_GROUP_MEMBERSHIP
	)
}
