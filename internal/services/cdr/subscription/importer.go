package subscription

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func importSubscriptionContext(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("expecting fhir_store,subscription_id,fhir_version as import string")
	}
	fhirStore := parts[0]
	id := parts[1]
	version := parts[2]

	if !slices.Contains([]string{"stu3", "r4"}, version) {
		return nil, fmt.Errorf("unsupported FHIR version '%s', myst be 'stu3' or 'r4'", version)
	}

	d.SetId(id)
	_ = d.Set("version", version)
	_ = d.Set("fhir_store", fhirStore)
	return []*schema.ResourceData{d}, nil
}
