package org_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	stu3pb "github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestParameters(t *testing.T) {
	body := `{
    "resourceType": "Parameters",
    "parameter": [
        {
            "name": "status",
            "valueString": "SUCCESS"
        },
        {
            "name": "submissionTime",
            "valueDateTime": "2018-07-13T00:05:29.981681+00:00"
        },
        {
            "name": "lastUpdated",
            "valueDateTime": "2018-07-13T00:08:37.563+00:00"
        },
        {
            "name": "requestor",
            "valueString": "77d6e95d-6f2a-4739-9d9c-bfa52f39a3e9"
        }
    ]
}`
	um, err := jsonformat.NewUnmarshaller("UTC", fhirversion.STU3)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, um) {
		return
	}
	unmarshalled, err := um.Unmarshal([]byte(body))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, unmarshalled) {
		return
	}
	contained := unmarshalled.(*stu3pb.ContainedResource)
	params := contained.GetParameters()

	assert.Len(t, params.Parameter, 4)
	assert.Equal(t, "status", params.Parameter[0].Name.Value)
	assert.Equal(t, "SUCCESS", params.Parameter[0].Value.GetStringValue().Value)
}

func TestAccResourceCDROrg_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_cdr_org.test"
	parentOrgID := acc.AccIAMOrgGUID()
	cdrURL := acc.AccCDRURL()

	if cdrURL == "" {
		t.Skipped()
		return
	}

	randomNameSTU3 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randomNameR4 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	now := time.Now().Format(time.RFC3339)

	// STU3
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceCDROrg(cdrURL, parentOrgID, randomNameSTU3, "stu3", now),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("MAIN %s %s", randomNameSTU3, now)),
				),
			},
		},
	})

	// R4
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceCDROrg(cdrURL, parentOrgID, randomNameR4, "r4", now),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("MAIN %s %s", randomNameR4, now)),
				),
			},
		},
	})
}

func testAccResourceCDROrg(cdrURL, parentOrgID, name, version, now string) string {
	return fmt.Sprintf(`

data "hsdp_cdr_fhir_store" "sandbox" {
  base_url = "%s"
  fhir_org_id = hsdp_iam_org.test.id
}

data "hsdp_iam_introspect" "terraform" {
}

data "hsdp_iam_service" "service" {
  service_id = data.hsdp_iam_introspect.terraform.subject
}

resource "hsdp_iam_org" "test" {
  name  = "MAIN%s"
  description = "CDR Org Acceptance Test %s %s"
  parent_org_id = "%s"
}

resource "hsdp_iam_org" "suborg" {
  name  = "SUB%s"
  description = "CDR Org Acceptance Test %s %s"
  parent_org_id = hsdp_iam_org.test.id
}

resource "hsdp_iam_role" "cdr_admin" {
  managing_organization = hsdp_iam_org.test.id
  name                  = "TF_CDR_ADMIN"
  permissions = [
    "ALL.READ",
    "ALL.WRITE"
  ]
}

resource "hsdp_iam_group" "cdr_admins" {
  managing_organization = hsdp_iam_org.test.id
  name                  = "TF_CDR_ADMIN"
  description           = "CDR Admins %s"
  roles                 = [hsdp_iam_role.cdr_admin.id]
  users                 = []
  services              = [data.hsdp_iam_service.service.id]
}

resource "hsdp_cdr_org" "test" {
  fhir_store  = data.hsdp_cdr_fhir_store.sandbox.endpoint

  name        = "MAIN %s %s"
  org_id      = hsdp_iam_org.test.id

  version     = "%s"

  depends_on = [hsdp_iam_group.cdr_admins]
}

resource "hsdp_cdr_org" "suborg" {
  fhir_store  = data.hsdp_cdr_fhir_store.sandbox.endpoint

  name        = "SUB %s %s"
  org_id      = hsdp_iam_org.suborg.id

  version     = "%s"
  
  part_of     = hsdp_cdr_org.test.org_id
}
`,
		// DATA SOURCE
		cdrURL,

		// IAM ORG
		name,
		now,
		name,
		parentOrgID,
		// IAM SUB-ORG
		name,
		now,
		name,

		// IAM GROUP
		now,

		// CDR ORG
		name,
		now,
		version,

		// CDR SUB ORG
		name,
		now,
		version,
	)
}
