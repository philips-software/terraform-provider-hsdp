package application_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMApplication_basic(t *testing.T) {
	t.Parallel()

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	parentOrgID := acc.AccIAMOrgGUID()

	upperRandomName := strings.ToUpper(randomName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "hsdp_iam_application.test",
				Config:       testAccResourceIAMApplication(parentOrgID, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hsdp_iam_application.test", "name", "ACC-"+upperRandomName),
				),
			},
		},
	})
}

func testAccResourceIAMApplication(parentOrgID, name string) string {
	// We create a completely separate ORG as that is currently
	// the only way we can clean up Propositions and Applications
	// after we are done testing
	upperName := strings.ToUpper(name)

	return fmt.Sprintf(`

resource "hsdp_iam_org" "test" {
  name = "ACC-%s"
  description = "IAM Application Test %s"

  parent_org_id = "%s"
  wait_for_delete = true
}

resource "hsdp_iam_proposition" "test" {
   name = "ACC-%s"
   description = "IAM Application Test %s"
   
   organization_id = hsdp_iam_org.test.id
}

resource "hsdp_iam_application" "test" {
    name = "ACC-%s"
    description = "IAM Application Test %s"
    proposition_id = hsdp_iam_proposition.test.id

    wait_for_delete = true
}`,
		// ORG
		upperName,
		name,
		parentOrgID,
		// PROP
		upperName,
		name,
		// APP
		upperName,
		name)
}
