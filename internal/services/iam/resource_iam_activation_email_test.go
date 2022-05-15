package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acctest"
)

func TestAccResourceIAMActivationEmail_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_iam_activation_email.test"
	userID := acctest.AccUserGUID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIAMActivationEmailConfig(userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_id", userID),
				),
			},
		},
	})
}

func testAccResourceIAMActivationEmailConfig(id string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_activation_email" "test" {
	user_id = %[1]q
}`, id)
}
