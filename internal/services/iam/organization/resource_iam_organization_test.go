package organization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMOrganization_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_org.test"
	parentOrgID := acc.AccIAMOrgGUID()
	randomOrgName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMOrganization(parentOrgID, randomOrgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "parent_org_id", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMOrganization(parentOrgID, name string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
    name = "ACCTest-%s"
    description = "ACC Test Org %s"

	parent_org_id = "%s"
    wait_for_delete = false
}`, name, name, parentOrgID)
}
