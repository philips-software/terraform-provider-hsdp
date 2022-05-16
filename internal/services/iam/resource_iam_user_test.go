package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMUser_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_user.test"
	parentOrgID := acc.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMUser(parentOrgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "organization_id", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMUser(parentOrgID string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_user" "test" {
  login           = "developer"
  email           = "acceptance@terrakube.com"
  first_name      = "Devel"
  last_name       = "Oper"
  password        = "DoNot@123"
  organization_id = %[1]q
 
}`, parentOrgID)
}
