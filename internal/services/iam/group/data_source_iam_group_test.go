package group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDataSourceIAMGroup_basic(t *testing.T) {

	org := acc.AccIAMOrgGUID()
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "hsdp_iam_group.test",
				Config:       testAccDataSourceIAMGroupORG(org, name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hsdp_iam_group.test", "name", name),
				),
			},
			{
				ResourceName: "data.hsdp_iam_group.test",
				Config:       testAccDataSourceIAMGroup(org, name, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hsdp_iam_group.test", "name", name),
					resource.TestCheckResourceAttr("data.hsdp_iam_group.test", "users.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceIAMGroup(org, name, randomName string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_user" "test" {
  login           = "%s"
  email           = "acceptance+%s@terrakube.com"
  first_name      = "ACC"
  last_name       = "Developer"
  password        = "DoNot@123"
  organization_id = "%s"
}

resource "hsdp_iam_group" "test" {
  name = "%s"
  description = "ACC Group DataSource Test %s"
  managing_organization = "%s"
  roles = []
  users = [hsdp_iam_user.test.id]
  services = []
}

data "hsdp_iam_group" "test" {
  managing_organization_id = "%s"
  name = "%s"

  depends_on = [hsdp_iam_group.test]
}`,
		// USER
		randomName,
		randomName,
		org,

		// RESOURCE
		name, name, org,
		// DATA
		org, name)
}

func testAccDataSourceIAMGroupORG(parentOrgID, name, randomName string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_group" "test" {
  name = "%s"
  description = "ACC Group DataSource Test %s"
  managing_organization = "%s"
  roles = []
  users = []
  services = []
}`, name, name, parentOrgID)
}
