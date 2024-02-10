package discovery_test

import (
	"fmt"
	"testing"

	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDatasourceDiscoveryService_basic(t *testing.T) {
	t.Parallel()

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := "data.hsdp_discovery_service.test"
	parentOrgID := acc.AccMDMOrgID()
	clientID := acc.AccMDMClientID()
	clientSecret := acc.AccMDMClientSecret()
	randomPassword, _ := tools.RandomPassword()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccDatasourceDiscoveryService(parentOrgID, randomName, clientID, clientSecret, randomPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "TerraformAcceptanceT"),
				),
			},
		},
	})
}

func testAccDatasourceDiscoveryService(parentOrgID, name, clientID, clientSecret, randomPassword string) string {
	// We create a completely separate ORG as that is currently
	// the only way we can clean up Propositions and Applications
	// after we are done testing

	return fmt.Sprintf(`
data "hsdp_iam_introspect" "test" {
  principal {
    username             = hsdp_iam_user.test.login
    password             = hsdp_iam_user.test.password
	oauth2_client_id     = "%s"
    oauth2_password      = "%s"
  }
}

resource "hsdp_iam_user" "test" {
  login           = "%s"
  email           = "%s@terrakube.com"
  first_name      = "ACC"
  last_name       = "TestUser"
  password        = "%s"

  organization_id = "%s" # Use parentOrgID
}

data "hsdp_discovery_service" "test" {
  principal {
    username    = hsdp_iam_user.test.login
    password    = hsdp_iam_user.test.password
    oauth2_client_id     = "%s"
    oauth2_password      = "%s"
  }

  tag = "terraform-acceptance"
}
`,
		// INTROSPECT
		clientID,
		clientSecret,
		// IAM USER
		name,
		name,
		randomPassword,
		parentOrgID,
		// DISCOVERY SERVICE
		clientID,
		clientSecret)
}
