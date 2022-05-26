package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDataSourceIAMGroup_basic(t *testing.T) {

	org := acc.AccIAMOrgGUID()
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("data.hsdp_iam_group.%s", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: fmt.Sprintf("hsdp_iam_group.%s", name),
				Config:       testAccDataSourceIAMGroupORG(org, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("hsdp_iam_group.%s", name), "name", name),
				),
			},
			{
				ResourceName: resourceName,
				Config:       testAccDataSourceIAMGroup(org, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func testAccDataSourceIAMGroup(org, name string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_group" "%s" {
  name = "%s"
  description = "ACC Group DataSource Test %s"
  managing_organization = "%s"
  roles = []
  users = []
  services = []
}

data "hsdp_iam_group" "%s" {
  managing_organization_id = "%s"
  name = "%s"
}`,
		// RESOURCE
		name, name, name, org,
		// DATA
		name, org, name)
}

func testAccDataSourceIAMGroupORG(parentOrgID, name string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_group" "%s" {
  name = "%s"
  description = "ACC Group DataSource Test %s"
  managing_organization = "%s"
  roles = []
  users = []
  services = []
}`, name, name, name, parentOrgID)
}
