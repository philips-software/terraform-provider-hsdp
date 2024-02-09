package notification_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceNotificationProducer_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_notification_producer.principal_producer"
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	iamOrgID := acc.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNotificationProducer(randomName, iamOrgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("acc principal producer %s", randomName)),
				),
			},
			{
				Config: testAccResourceNotificationProducer(randomName, iamOrgID) +
					fmt.Sprintf(`
resource "hsdp_notification_producer" "producer" {
  managing_organization_id       = hsdp_iam_org.test.id
  managing_organization          = hsdp_iam_org.test.name
  producer_product_name          = "accProduct%s"
  producer_service_name          = "accService%s"
  producer_service_instance_name = "accServiceInstance%s"
  producer_service_base_url      = "https://ns-producer.terrakube.com/"
  producer_service_path_url      = "notification/create/%s"
  description                    = "acc producer %s"

  depends_on = [hsdp_iam_group.producer_admins]
}`, randomName, randomName, randomName, randomName, randomName),
				ExpectError: regexp.MustCompile("User is not in the given organization"),
			},
		},
	})
}

func testAccResourceNotificationProducer(random, parentId, password string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
    name = "ACCTest-%s"
    description = "ACCResourceNotificationProducer %s"

	parent_org_id = "%s"
    wait_for_delete = false
}

resource "hsdp_iam_proposition" "test" {
   name = "ACC-%s"
   description = "IAM Service Test %s"
   
   organization_id = hsdp_iam_org.test.id
}

resource "hsdp_iam_application" "test" {
    name = "ACC-%s"
    description = "IAM Service Test %s"
    proposition_id = hsdp_iam_proposition.test.id
}

resource "hsdp_iam_service" "test" {
  name        = "ACC-%s"
  description = "IAM Service Test %s"

  validity = 12
  token_validity = 3600
  scopes = ["openid"]
  default_scopes = ["openid"]
 
  application_id = hsdp_iam_application.test.id
}

resource "hsdp_iam_role" "producer_admin" {
  name = "PRODUCER_ADMIN_TF"
  permissions = [
    "ORGANIZATION.READ",
    "NS_PRODUCER.CREATE",
    "NS_PRODUCER.READ",
    "NS_PRODUCER.DELETE",
    "NS_TOPIC_SCOPE.READ",
    "NS_TOPIC.CREATE",
    "NS_TOPIC.READ",
    "NS_TOPIC.UPDATE",
    "NS_TOPIC.DELETE",
    "NS_SUBSCRIBER.READ",
    "NS_SUBSCRIPTION.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_role" "publisher" {
  name = "PUBLISHER_TF"
  permissions = [
    "NS_PUBLISH.CREATE",
    "NS_TOPIC.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_role" "subscriber_admin" {
  name = "SUBSCRIBER_ADMIN_TF"
  permissions = [
    "ORGANIZATION.READ",
    "NS_SUBSCRIBER.CREATE",
    "NS_SUBSCRIBER.DELETE",
    "NS_PRODUCER.READ",
    "NS_SUBSCRIPTION.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_role" "subscriber" {
  name = "SUBSCRIBER_TF"
  permissions = [
    "NS_SUBSCRIPTION.CREATE",
    "NS_SUBSCRIPTION.READ",
    "NS_SUBSCRIPTION.DELETE",
    "NS_SUBSCRIPTION.CONFIRM",
    "NS_SUBSCRIPTION.SYNC",
    "NS_TOPIC.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "producer_admins" {
  name                  = "PRODUCER_ADMINS_TF"
  roles                 = [hsdp_iam_role.producer_admin.id]
  services              = [hsdp_iam_service.test.id]
  users                 = [hsdp_iam_user.user.id]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "publishers" {
  name                  = "PUBLISHERS_TF"
  roles                 = [hsdp_iam_role.publisher.id]
  users                 = [hsdp_iam_user.user.id]
  services              = [hsdp_iam_service.test.id]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "subscriber_admins" {
  name                  = "SUBSCRIBER_ADMINS_TF"
  roles                 = [hsdp_iam_role.subscriber_admin.id]
  users                 = [hsdp_iam_user.user.id]
  services              = [hsdp_iam_service.test.id]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "subscribers" {
  name                  = "SUBSCRIBERS_TF"
  roles                 = [hsdp_iam_role.subscriber.id]
  users                 = [hsdp_iam_user.user.id]
  services              = [hsdp_iam_service.test.id]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_notification_producer" "principal_producer" {
  principal {
    service_id          = hsdp_iam_service.test.service_id
    service_private_key = hsdp_iam_service.test.private_key
  }

  managing_organization_id       = hsdp_iam_org.test.id
  managing_organization          = hsdp_iam_org.test.name
  producer_product_name          = "accPrincipalProduct%s"
  producer_service_name          = "accPrincipalService%s"
  producer_service_instance_name = "accPrincipalServiceInstance%s"
  producer_service_base_url      = "https://ns-principal-producer.terrakube.com/"
  producer_service_path_url      = "notification/create/%s"
  description                    = "acc principal producer %s"

  depends_on = [hsdp_iam_group.producer_admins]
}

resource "hsdp_iam_user" "user" { 
  organization_id = hsdp_iam_org.test.id
  login           = "login-%s"
  password        = "%s"
  email           = "user-%s@terrakube.com"
  first_name      = "ACCProducer"
  last_name       = "ACCAcount"
}
`,
		// IAM TEST ORG
		random,
		random,
		parentId,

		// IAM Proposition
		strings.ToUpper(random),
		random,

		// IAM Application
		strings.ToUpper(random),
		random,

		// IAM Service
		strings.ToUpper(random),
		random,

		// IAM GROUP producer_admins

		// IAM GROUP publishers

		// IAM GROUP subscriber_admins

		// IAM GROUP subscribers

		// NOTIFICATION PRINCIPAL PRODUCER producer
		random,
		random,
		random,
		random,
		random,

		// IAM USER
		random,
                password,
		random,
	)
}
