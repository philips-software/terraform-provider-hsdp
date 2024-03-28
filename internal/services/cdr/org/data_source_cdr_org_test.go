package org_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccDataSourceCDROrg_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.hsdp_cdr_org.test"
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
				Config:       testAccDataSourceCDROrg(cdrURL, parentOrgID, randomNameSTU3, "stu3", now),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s %s", randomNameSTU3, now)),
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
				Config:       testAccDataSourceCDROrg(cdrURL, parentOrgID, randomNameR4, "r4", now),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s %s", randomNameR4, now)),
				),
			},
		},
	})
}

func testAccDataSourceCDROrg(cdrURL, parentOrgID, name, version, now string) string {
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
  name  = "%s"
  description = "Data Source ORG Acceptance Test CDR %s %s"
  parent_org_id = "%s"
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

  name        = "%s %s"
  org_id      = hsdp_iam_org.test.id

  version     = "%s"
}

data "hsdp_cdr_org" "test" {
  fhir_store  = data.hsdp_cdr_fhir_store.sandbox.endpoint

  org_id      = hsdp_cdr_org.test.id

  version     = "%s"
}
`,
		// DATA SOURCE
		cdrURL,

		// IAM ORG
		name,
		now,
		name,
		parentOrgID,

		// IAM GROUP
		now,

		// CDR ORG
		name,
		now,
		version,

		// DATA CDR ORG
		version,
	)
}
