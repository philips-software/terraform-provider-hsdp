package group_test

import (
	"fmt"
	"strings"
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
					resource.TestCheckResourceAttr(resourceName, "devices.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "services.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
		},
	})
}

func testAccResourceIAMGroup(parentOrgID, name string) string {
	upperName := strings.ToUpper(name)

	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
  name = "ACC-%s"
  description = "IAM Group Test %s"

  parent_org_id = "%s"
  wait_for_delete = false
}

resource "hsdp_iam_proposition" "test" {
   name = "ACC-%s"
   description = "IAM Group Test %s"
   
   organization_id = hsdp_iam_org.test.id
}

resource "hsdp_iam_application" "test" {
    name = "ACC-%s"
    description = "IAM Group Test %s"
    proposition_id = hsdp_iam_proposition.test.id
}

resource "hsdp_iam_device" "test" {
  login_id = "d%s"
  password = "Random%s@123!"

  organization_id = hsdp_iam_org.test.id
  application_id  = hsdp_iam_application.test.id

  external_identifier {
    type {
      code = "ID"
      text = "Device Identifier"
    }
    system = "http://www.philips.co.id/phs/healthwatch"
    value = "%s"
  }

  type = "ActivityMonitor"

  for_test = true
  is_active = true
}

resource "hsdp_iam_service" "test" {
  name        = "ACC-%s"
  description = "IAM Group Test %s"

  validity = 12
  token_validity = 3600
  scopes = ["openid"]
  default_scopes = ["openid"]
 
  application_id = hsdp_iam_application.test.id
}

resource "hsdp_iam_user" "test" {
  login           = "u%s"
  email           = "acceptance+%s@terrakube.com"
  first_name      = "ACC"
  last_name       = "Developer"
  password        = "DoNot@123"
  organization_id = hsdp_iam_org.test.id
 
}

resource "hsdp_iam_role" "test" {
  name        = "%s"
  description = "Acceptance Test Group %s"

  permissions = [
    "ALL.READ",
    "ALL.WRITE"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "test" {
  name        = "%s"
  description = "Acceptance Test Group %s"

  roles    = [hsdp_iam_role.test.id]
  users    = [hsdp_iam_user.test.id]
  services = [hsdp_iam_service.test.id]
  devices  = [hsdp_iam_device.test.id]

  managing_organization = hsdp_iam_org.test.id
 
  drift_detection = true
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
		name,
		// DEVICE
		name,
		name,
		name,
		// SERVICE,
		upperName,
		name,
		// USER
		name,
		name,
		// ROLE
		upperName,
		name,
		// GROUP
		name,
		name)
}
