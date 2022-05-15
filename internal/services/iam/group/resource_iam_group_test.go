package group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acctest"
)

func TestAccResourceIAMGroup_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_group.test"
	parentOrgID := acctest.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMGroup(parentOrgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "managing_organization", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMGroup(parentOrgID string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_group" "test" {
  name        = "TESTGROUP"
  description = "Acceptance Test Group"

  roles    = []
  users    = []
  services = []

  managing_organization = %[1]q
 
  drift_detection = true
}`, parentOrgID)
}
