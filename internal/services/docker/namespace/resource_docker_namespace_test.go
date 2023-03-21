package namespace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceDockerNamespace_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_docker_namespace.test"
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceDockerNamespace(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randomName),
				),
			},
		},
	})
}

func testAccResourceDockerNamespace(name string) string {
	return fmt.Sprintf(`
resource "hsdp_docker_namespace" "test" {
  name = "%s"
}
`, name)
}
