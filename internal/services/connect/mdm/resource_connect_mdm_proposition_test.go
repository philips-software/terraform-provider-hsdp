package mdm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceMDMProposition_basic(t *testing.T) {
	t.Parallel()

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("hsdp_connect_mdm_proposition.%s", randomName)
	parentOrgID := acc.AccIAMOrgGUID()

	upperRandomName := strings.ToUpper(randomName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceMDMProposition(parentOrgID, randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ACC-"+upperRandomName),
				),
			},
		},
	})
}

func testAccResourceMDMProposition(parentOrgID, name string) string {
	// We create a completely separate ORG as that is currently
	// the only way we can clean up Propositions and Applications
	// after we are done testing
	upperName := strings.ToUpper(name)

	return fmt.Sprintf(`

data "hsdp_iam_introspect" "%s" {
}

data "hsdp_iam_service" "%s" {
  service_id = data.hsdp_iam_introspect.%s.subject
}

resource "hsdp_iam_org" "%s" {
  name = "ACC-%s"
  description = "MDM Proposition Test %s"

  parent_org_id = "%s"
  wait_for_delete = false
}

resource "hsdp_iam_role" "test" {
  name = "ACC-%s"
  description = "MDM Proposition Test %s"
  permissions = [
    # Set 1
    # Read only access is needed to query global MDM resources.
    "MDM-REGION.READ",
    "MDM-STANDARDSERVICE.READ",
    "MDM-STORAGECLASS.READ",
    "MDM-OAUTHCLIENTSCOPE.READ",
    "MDM-SUBSCRIBERTYPE.READ",
    "MDM-DATASUBSCRIBER.READ",
    "MDM-SERVICEAGENT.READ",
    "MDM-AUTHENTICATIONMETHOD.READ",
    "MDM-DATAADAPTER.READ",
    "MDM-SERVICEACTION.READ",

    # Set 2
    # For managing your core hierarchy using APIs, assign these permissions.
    "PROPOSITION.WRITE",
    "MDM-PROPOSITION.CREATE",
    "MDM-PROPOSITION.READ",
    "MDM-PROPOSITION.UPDATE",

    "APPLICATION.WRITE",
    "MDM-APPLICATION.CREATE",
    "MDM-APPLICATION.READ",
    "MDM-APPLICATION.UPDATE",

    "GROUP.READ",

    "MDM-DEVICEGROUP.CREATE",
    "MDM-DEVICEGROUP.READ",
    "MDM-DEVICEGROUP.UPDATE",
    "MDM-DEVICEGROUP.DELETE",

    "MDM-DEVICETYPE.CREATE",
    "MDM-DEVICETYPE.READ",
    "MDM-DEVICETYPE.UPDATE",
    "MDM-DEVICETYPE.DELETE",

    "CLIENT.WRITE",
    "CLIENT.READ",
    # CLIENT.SCOPES
    "CLIENT.DELETE",
    "MDM-OAUTHCLIENT.CREATE",
    "MDM-OAUTHCLIENT.READ",
    "MDM-OAUTHCLIENT.UPDATE",
    "MDM-OAUTHCLIENT.DELETE",

    # Set 3
    # For managing discovery service master data using APIs, assign these permissions.
    "MDM-SERVICEREFERENCE.CREATE",
    "MDM-SERVICEREFERENCE.READ",
    "MDM-SERVICEREFERENCE.UPDATE",
    "MDM-SERVICEREFERENCE.DELETE",
    "MDM-SERVICEACTION.CREATE",
    "MDM-SERVICEACTION.READ",
    "MDM-SERVICEACTION.UPDATE",
    "MDM-SERVICEACTION.DELETE",

    # Set 4
    # For managing of firmware service master data using APIs, assign these permissions.
    "MDM-FIRMWARECOMPONENT.CREATE",
    "MDM-FIRMWARECOMPONENT.READ",
    "MDM-FIRMWARECOMPONENT.UPDATE",
    "MDM-FIRMWARECOMPONENT.DELETE",
    "MDM-FIRMWARECOMPONENTVERSION.CREATE",
    "MDM-FIRMWARECOMPONENTVERSION.READ",
    "MDM-FIRMWARECOMPONENTVERSION.UPDATE",
    "MDM-FIRMWARECOMPONENTVERSION.DELETE",
    "MDM-FIRMWAREDISTRIBUTIONREQUEST.CREATE",
    "MDM-FIRMWAREDISTRIBUTIONREQUEST.READ",
    "MDM-FIRMWAREDISTRIBUTIONREQUEST.UPDATE",

    # Set 5
    # For managing blob repository service master data using APIs, assign these permissions.
    "MDM-BUCKET.CREATE",
    "MDM-BUCKET.READ",
    "MDM-BUCKET.UPDATE",
    "MDM-BUCKET.DELETE",
    "MDM-DATATYPE.CREATE",
    "MDM-DATATYPE.READ",
    "MDM-DATATYPE.UPDATE",
    "MDM-DATATYPE.DELETE",
    "MDM-BLOBDATACONTRACT.CREATE",
    "MDM-BLOBDATACONTRACT.READ",
    "MDM-BLOBDATACONTRACT.UPDATE",
    "MDM-BLOBDATACONTRACT.DELETE",

    "MDM-BLOBSUBSCRIPTION.CREATE",
    "MDM-BLOBSUBSCRIPTION.READ",
    "MDM-BLOBSUBSCRIPTION.UPDATE",
    "MDM-BLOBSUBSCRIPTION.DELETE",

    # Set 6
    # For managing data broker service master data using APIs, assign these permissions.
    "MDM-DATABROKERSUBSCRIPTION.CREATE",
    "MDM-DATABROKERSUBSCRIPTION.READ",
    "MDM-DATABROKERSUBSCRIPTION.UPDATE",
    "MDM-DATABROKERSUBSCRIPTION.DELETE",

    "MDM-DATATYPE.CREATE",
    "MDM-DATATYPE.READ",
    "MDM-DATATYPE.UPDATE",
    "MDM-DATATYPE.DELETE",

    # Other
    "NS_TOPIC.READ",
    "NS_PRODUCER.READ",
    "MDM-AUTHENTICATIONMETHOD.CREATE"
  ]
  managing_organization = hsdp_iam_org.%s.id
}

resource "hsdp_iam_group" "test" {
  name = "ACC-%s"
  description = "MDM Proposition Test %s"
  roles = [hsdp_iam_role.test.id]
  services = [data.hsdp_iam_service.%s.id]
  managing_organization = hsdp_iam_org.%s.id
}

resource "hsdp_connect_mdm_proposition" "%s" {
   name = "ACC-%s"
   description = "MDM Proposition Test %s"
   
   organization_id = hsdp_iam_org.%s.id
   
   status = "ACTIVE"
}`,
		// INTROSPECT
		name,
		// DATA SOURCE
		name,
		name,
		// ORG
		name,
		upperName,
		name,
		parentOrgID,
		// ROLE
		upperName,
		name,
		name,
		// GROUP
		name,
		name,
		name,
		name,
		// PROP
		name,
		upperName,
		name,
		name)
}
