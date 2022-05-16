package group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMGroup_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_group.test"
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
				Config:       testAccResourceIAMGroup(parentOrgID, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "managing_organization", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMGroup(parentOrgID, name string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_group" "test" {
  name        = "%s"
  description = "Acceptance Test Group %s"

  roles    = []
  users    = []
  services = []

  managing_organization = "%s"
 
  drift_detection = true
}`, name, name, parentOrgID)
}
