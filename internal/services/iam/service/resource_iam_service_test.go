package service_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMService_basic(t *testing.T) {
	t.Parallel()

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("hsdp_iam_service.%s", randomName)
	parentOrgID := acc.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMService(parentOrgID, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "validity", "12"),
					resource.TestCheckResourceAttr(resourceName, "token_validity", "3600"),
				),
			},
		},
	})
}

func testAccResourceIAMService(parentOrgID, name string) string {
	// We create a completely separate ORG as that is currently
	// the only way we can clean up Propositions and Applications
	// after we are done testing
	upperName := strings.ToUpper(name)

	return fmt.Sprintf(`

resource "hsdp_iam_org" "%s" {
  name = "ACC-%s"
  description = "IAM Service Test %s"

  parent_org_id = "%s"
  wait_for_delete = false
}

resource "hsdp_iam_proposition" "%s" {
   name = "ACC-%s"
   description = "IAM Service Test %s"
   
   organization_id = hsdp_iam_org.%s.id
}

resource "hsdp_iam_application" "%s" {
    name = "ACC-%s"
    description = "IAM Service Test %s"
    proposition_id = hsdp_iam_proposition.%s.id
}

resource "hsdp_iam_service" "%s" {
  name        = "ACC-%s"
  description = "IAM Service Test %s"

  validity = 12
  token_validity = 3600
  scopes = ["openid"]
  default_scopes = ["openid"]
 
  application_id = hsdp_iam_application.%s.id
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
		name,
		// APP
		name,
		upperName,
		name,
		name,
		// SERVICE
		name,
		upperName,
		name,
		name)
}
