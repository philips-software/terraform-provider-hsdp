package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMUser_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_user.test"
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
				Config:       testAccResourceIAMUser(parentOrgID, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "organization_id", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMUser(parentOrgID, name string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_user" "test" {
  login           = "%s"
  email           = "acceptance+%s@terrakube.com"
  first_name      = "ACC"
  last_name       = "Developer"
  password        = "DoNot@123"
  organization_id = "%s"
 
}`, name, name, parentOrgID)
}
