package org

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	r4dt "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/organization_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func r4Create(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	orgID := d.Get("org_id").(string)
	name := d.Get("name").(string)
	partOf := d.Get("part_of").(string)
	org, err := r4.NewOrganization(c.TimeZone, orgID, name)
	if err != nil {
		return diag.FromErr(err)
	}
	if partOf != "" {
		org.PartOf = &r4dt.Reference{
			Reference: &r4dt.Reference_OrganizationId{
				OrganizationId: &r4dt.ReferenceId{
					Value: partOf,
				},
			},
		}
	}
	var onboardedOrg *r4pb.Organization

	err = tools.TryHTTPCall(ctx, 9, func() (*http.Response, error) {
		var resp *cdr.Response
		var err error
		onboardedOrg, resp, err = client.TenantR4.Onboard(org)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	}, append(tools.StandardRetryOnCodes, http.StatusNotFound)...) // CDR weirdness

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(onboardedOrg.Id.GetValue())
	return diags
}

func r4Read(ctx context.Context, client *cdr.Client, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	var org *r4pb.Organization
	var resp *cdr.Response

	id := d.Id()

	err := tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error

		org, resp, err = client.TenantR4.GetOrganizationByID(id)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	}, append(tools.StandardRetryOnCodes, http.StatusNotFound)...) // CDR weirdness

	if err != nil {
		if resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone) {
			d.SetId("")
			return diags
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	if org == nil {
		return diag.FromErr(fmt.Errorf("org is nil, this is unexpected (id='%s'): %w", id, err))
	}
	if org.Name != nil {
		_ = d.Set("name", org.Name.GetValue())
	}
	if org.PartOf != nil {
		partOfOrgID := org.PartOf.GetOrganizationId()
		_ = d.Set("part_of", partOfOrgID.GetValue())
	}
	return diags
}

func r4Update(_ context.Context, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)

	id := d.Id()
	org, _, err := client.TenantR4.GetOrganizationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonOrg, err := c.R4MA.MarshalResource(org)
	if err != nil {
		return diag.FromErr(err)
	}
	madeChanges := false

	if d.HasChange("name") {
		org.Name.Value = d.Get("name").(string)
		madeChanges = true
	}
	if d.HasChange("part_of") {
		partOf := d.Get("part_of").(string)
		if partOf != "" {
			org.PartOf = &r4dt.Reference{
				Reference: &r4dt.Reference_OrganizationId{
					OrganizationId: &r4dt.ReferenceId{
						Value: partOf,
					},
				},
			}
		} else {
			org.PartOf = nil
		}
		madeChanges = true
	}
	if !madeChanges {
		return diags
	}

	changedOrg, _ := c.R4MA.MarshalResource(org)
	patch, err := jsonpatch.DiffBytes(jsonOrg, changedOrg)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = client.OperationsR4.Patch("Organization/"+id, patch)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func r4PurgeStateRefreshFunc(client *cdr.Client, purgeStatusURL, id string) resource.StateRefreshFunc {
	statusURL, err := url.Parse(purgeStatusURL)
	if err != nil {
		return func() (result interface{}, state string, err error) {
			return nil, "FAILED", err
		}
	}
	return func() (interface{}, string, error) {
		contained, resp, err := client.OperationsR4.Get(id, func(request *http.Request) error {
			request.URL = statusURL
			return nil
		})
		if err != nil {
			return resp, "FAILED", err
		}
		if resp.StatusCode() == http.StatusAccepted { // In progress
			return resp, "PURGING", nil
		}
		params := contained.GetParameters()
		// Return the status value
		for _, p := range params.Parameter {
			if p.Name.Value == "status" {
				return resp, p.Value.GetStringValue().Value, nil
			}
		}
		return resp, "FAILED", fmt.Errorf("missing status parameter for GET %s", purgeStatusURL)
	}
}

func r4Delete(ctx context.Context, client *cdr.Client, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Id()
	purgeDelete := d.Get("purge_delete").(bool)

	if !purgeDelete {
		deleted, resp, err := client.OperationsR4.Delete(path.Join("Organization", id))
		if resp != nil && resp.StatusCode() == http.StatusNotFound { // Already gone
			d.SetId("")
			return diags
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if !deleted {
			if resp != nil {
				return diag.FromErr(fmt.Errorf("delete failed with status code %d", resp.StatusCode()))
			}
			return diag.FromErr(fmt.Errorf("delete failed with nil response"))
		}
		d.SetId("")
		return diags
	}
	// Purge delete with purge-status check
	_, resp, err := client.OperationsR4.Post(path.Join("$purge"), []byte(``), func(request *http.Request) error {
		request.URL.Opaque = "/store/fhir/" + id + "/$purge"
		return nil
	})
	if resp != nil && resp.StatusCode() == http.StatusNotFound { // Already gone
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.FromErr(fmt.Errorf("unexpected nil response for $purge operation"))
	}
	if resp.StatusCode() != http.StatusAccepted {
		return diag.FromErr(fmt.Errorf("$purge operation returned unexpected statusCode %d", resp.StatusCode()))
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PURGING"},
		Target:     []string{"SUCCESS"},
		Refresh:    r4PurgeStateRefreshFunc(client, resp.Header.Get("Location"), id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"error waiting for FHIR ORG purge '%s' operation: %v",
			id, err))
	}
	d.SetId("")
	return diags
}
