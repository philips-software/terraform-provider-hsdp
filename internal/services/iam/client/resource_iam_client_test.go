package client_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func TestAccResourceIAMClient_basic(t *testing.T) {
	t.Parallel()

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := "hsdp_iam_client.test"
	parentOrgID := acc.AccIAMOrgGUID()
	description1 := "ACCFOO"
	description2 := "ACCBAR"
	name1 := fmt.Sprintf("%s-1", randomName)
	name2 := fmt.Sprintf("%s-2", randomName)
	randomPassword, _ := tools.RandomPassword()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMClient(parentOrgID, "Confidential", randomName, name1, description1, randomPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "client_id", "acc-"+randomName),
					resource.TestCheckResourceAttr(resourceName, "description", description1),
				),
			},
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMClient(parentOrgID, "Confidential", randomName, name2, description2, randomPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "client_id", "acc-"+randomName),
					resource.TestCheckResourceAttr(resourceName, "description", description2),
					resource.TestCheckResourceAttr(resourceName, "type", "Confidential"),
				),
			},
		},
	})
}

func testAccResourceIAMClient(parentOrgID, clientType, fixedName, changingName, description, randomPassword string) string {
	// We create a completely separate ORG as that is currently
	// the only way we can clean up Propositions and Applications
	// after we are done testing
	upperName := strings.ToUpper(fixedName)
	clientId := "acc-" + fixedName

	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
  name = "ACC-%s"
  description = "IAM Application Test %s"

  parent_org_id = "%s"
  wait_for_delete = false
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
}

resource "hsdp_iam_client" "test" {
    type                = "%s"
    name                = "%s"
    client_id           = "%s"
    password            = "%s"
    application_id      = hsdp_iam_application.test.id
    global_reference_id = "678477ff-35cb-4999-9100-0e74a16b820b"
    description         = "%s"

    scopes = [ "cn", "auth_iam_introspect", "email", "profile" ]

    default_scopes = [ "cn", "auth_iam_introspect" ]

    redirection_uris = [
      "https://foo.bar/auth",
      "https://testapp.cloud.pcftest.com/auth",
    ]

    response_types = ["code", "code id_token"]
}
`,
		// ORG
		upperName,
		fixedName,
		parentOrgID,
		// PROP
		upperName,
		fixedName,
		// APP
		upperName,
		fixedName,
		// CLIENT
		clientType,
		changingName,
		clientId,
		randomPassword,
		description,
	)
}
