package proposition_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMProposition_basic(t *testing.T) {
	t.Parallel()

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("hsdp_iam_proposition.%s", randomName)
	parentOrgID := acc.AccIAMOrgGUID()

	upperRandomName := strings.ToUpper(randomName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMProposition(parentOrgID, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ACC-"+upperRandomName),
				),
			},
		},
	})
}

func testAccResourceIAMProposition(parentOrgID, name string) string {
	// We create a completely separate ORG as that is currently
	// the only way we can clean up Propositions and Applications
	// after we are done testing
	upperName := strings.ToUpper(name)

	return fmt.Sprintf(`

resource "hsdp_iam_org" "%s" {
  name = "ACC-%s"
  description = "IAM Proposition Test %s"

  parent_org_id = "%s"
  wait_for_delete = false
}

resource "hsdp_iam_proposition" "%s" {
   name = "ACC-%s"
   description = "IAM Proposition Test %s"
   
   organization_id = hsdp_iam_org.%s.id
}`,
		// ORG
		name,
		upperName,
		name,
		parentOrgID,
		// PROP
		name,
		upperName,
		name,
		name)
}
