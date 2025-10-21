package iot_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceConnectIoTProvisioningOrgConfiguration_basic(t *testing.T) {
	t.Parallel()
	resourceName := "hsdp_connect_iot_provisioning_orgconfiguration.test"
	orgGuid := acc.AccIAMOrgGUID()
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceConnectIoTProvisioningOrgConfiguration(orgGuid, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "organization_guid", orgGuid),
					resource.TestCheckResourceAttr(resourceName, "service_account.0.service_account_id", fmt.Sprintf("demo_test_tf_%s.app@demo.iot__connect__sandbox.apmplatform.philips-healthsuite.com", randomName)),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_signature.0.algorithm", "RSA-SHA256"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_signature.0.config.0.type", "RSA"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_signature.0.config.0.padding", "RSA_PKCS1_PSS_PADDING"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_signature.0.config.0.salt_length", "RSA_PSS_SALTLEN_MAX_SIGN"),
				),
			},
		},
	})
}

func testAccResourceConnectIoTProvisioningOrgConfiguration(orgGuid, randomName string) string {
	return fmt.Sprintf(`
resource "hsdp_connect_iot_provisioning_orgconfiguration" "test" {
  organization_guid = "%s"
  
  service_account {
    service_account_id  = "demo_test_tf_%s.app@demo.iot__connect__sandbox.apmplatform.philips-healthsuite.com"
    service_account_key = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAwJ6bR9Wj3wSxHGIBxmO8VVLhLUGzXGJFVdgzjJwMIIEowIBAAKCAQEAwJ6bR9Wj3wSxHGIBxmO8VVLhLUGzXGJFVdgzjJwMKExUm\n-----END RSA PRIVATE KEY-----"
  }

  bootstrap_signature {
    algorithm  = "RSA-SHA256"
    public_key = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCoALFgprtmwkm7jF5kqZmF3XmVlHRjF6rQWqMqgzQIDAQAB\n-----END PUBLIC KEY-----"
    
    config {
      type        = "RSA"
      padding     = "RSA_PKCS1_PSS_PADDING"
      salt_length = "RSA_PSS_SALTLEN_MAX_SIGN"
    }
  }
}`, orgGuid, randomName)
}