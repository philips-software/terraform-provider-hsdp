package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acctest"
)

func TestAccResourceIAMOrganization_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_org.test"
	parentOrgID := acctest.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMOrganization(parentOrgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "parent_org_id", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMOrganization(parentOrgID string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
    name = "ACCTestORG"
    description = "ACCTestORG"
	parent_org_id = %[1]q
    wait_for_delete = true
}`, parentOrgID)
}
