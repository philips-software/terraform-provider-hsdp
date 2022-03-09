package org

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/cenkalti/backoff/v4"
	r4dt "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/organization_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func r4Create(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	onboardedOrg, _, err = client.TenantR4.GetOrganizationByID(orgID)
	if err == nil && onboardedOrg != nil {
		d.SetId(onboardedOrg.Id.Value)
		return resourceCDROrgUpdate(ctx, d, m)
	}
	// Do initial boarding
	operation := func() error {
		var resp *cdr.Response
		onboardedOrg, resp, err = client.TenantR4.Onboard(org)
		if resp == nil {
			if err != nil {
				return err
			}
			return fmt.Errorf("TenantR4.Onboard: response is nil")
		}
		// TODO: refactor this check so we don't have to check for nil resp above
		return tools.CheckForIAMPermissionErrors(client, resp.Response, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(onboardedOrg.Id.Value)
	return diags
}

func r4Read(_ context.Context, client *cdr.Client, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	orgID := d.Get("org_id").(string)
	org, resp, err := client.TenantR4.GetOrganizationByID(orgID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
		return diags
	}
	_ = d.Set("name", org.Name.Value)
	if org.PartOf != nil {
		_ = d.Set("part_of", org.PartOf.GetOrganizationId())
	}
	return diags
}

func r4Update(ctx context.Context, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		if resp.StatusCode == http.StatusAccepted { // In progress
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

func r4Delete(ctx context.Context, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Id()
	purgeDelete := d.Get("purge_delete").(bool)

	if !purgeDelete {
		deleted, resp, err := client.OperationsR4.Delete(path.Join("Organization", id))
		if resp != nil && resp.StatusCode == http.StatusNotFound { // Already gone
			d.SetId("")
			return diags
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if !deleted {
			if resp != nil {
				return diag.FromErr(fmt.Errorf("delete failed with status code %d", resp.StatusCode))
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
	if resp != nil && resp.StatusCode == http.StatusNotFound { // Already gone
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.FromErr(fmt.Errorf("unexpected nil response for $purge operation"))
	}
	if resp.StatusCode != http.StatusAccepted {
		return diag.FromErr(fmt.Errorf("$purge operation returned unexpected statusCode %d", resp.StatusCode))
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
