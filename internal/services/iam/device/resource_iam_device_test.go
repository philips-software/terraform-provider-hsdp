package device_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func TestAccResourceIAMDevice_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_device.test"
	parentOrgID := acc.AccIAMOrgGUID()
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randomPassword, _ := tools.RandomPassword()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceIAMUser(parentOrgID, randomName, randomPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "login_id", randomName),
				),
			},
		},
	})
}

func testAccResourceIAMUser(parentOrgID, name, randomPassword string) string {
	upperName := strings.ToUpper(name)

	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
  name = "ACC-%s"
  description = "IAM Device Test %s"

  parent_org_id = "%s"
  wait_for_delete = false
}

resource "hsdp_iam_proposition" "test" {
   name = "ACC-%s"
   description = "IAM Device Test %s"
   
   organization_id = hsdp_iam_org.test.id
}

resource "hsdp_iam_application" "test" {
    name = "ACC-%s"
    description = "IAM Device Test %s"
    proposition_id = hsdp_iam_proposition.test.id
}

resource "hsdp_iam_device" "test" {
  login_id = "%s"
  password = "%s"

  organization_id = hsdp_iam_org.test.id
  application_id  = hsdp_iam_application.test.id

  external_identifier {
    type {
      code = "ID"
      text = "Device Identifier"
    }
    system = "https://www.philips.co.id/phs/healthwatch"
    value = "%s"
  }

  type = "ActivityMonitor"

  for_test = true
  is_active = true
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
		randomPassword,
		name)
}
