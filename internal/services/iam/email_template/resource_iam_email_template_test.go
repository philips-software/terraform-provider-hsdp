package email_template_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceIAMEmailTemplate_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_email_template.account_already_verified"
	parentOrgID := acc.AccIAMOrgGUID()
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIAMEmailTemplateConfig(randomName, parentOrgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "managing_organization", parentOrgID),
				),
			},
		},
	})
}

func testAccResourceIAMEmailTemplateConfig(id, parentOrgId string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
	name = "ACCTest-%s"
    description = "ACC Email Template Test"

	parent_org_id = "%s"
    wait_for_delete = false
}
 
resource "hsdp_iam_email_template" "account_already_verified" {
    type = "ACCOUNT_ALREADY_VERIFIED"

    managing_organization = "%s"

    format   = "HTML"
    message = "Yo dawg, your account is already verified"
}

resource "hsdp_iam_email_template" "account_unlocked" {
    type = "ACCOUNT_UNLOCKED"

    managing_organization = "%s"

    format   = "HTML"
    message = "Yo dawg, I unlocked your account"
}

resource "hsdp_iam_email_template" "account_verification" {
    type = "ACCOUNT_VERIFICATION"

    managing_organization = "%s"

    format   = "HTML"
    message = "Yo dawg, verify your account here {{link.verification}} -- You have {{template.linkExpiryPeriod}} hours!"
}

resource "hsdp_iam_email_template" "password_expiry" {
    type = "PASSWORD_EXPIRY"

    managing_organization = "%s"
   
    locale = "en-US"
	from   = "ron.swanson@pawnee.city"

    format   = "HTML"
    message = "Yo dawg, your password is about to expire, change it!"
}

`,

		// ORG
		id,
		parentOrgId,

		// TEMPLATE ACCOUNT_ALREADY_VERIFIED
		parentOrgId,

		// TEMPLATE ACCOUNT_UNLOCKED
		parentOrgId,

		// TEMPLATE ACCOUNT_VERIFICATION
		parentOrgId,

		// TEMPLATE PASSWORD_EXPIRY
		parentOrgId,
	)
}
