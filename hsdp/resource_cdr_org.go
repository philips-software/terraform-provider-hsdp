package hsdp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceCDROrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDROrgCreate,
		ReadContext:   resourceCDROrgRead,
		UpdateContext: resourceCDROrgUpdate,
		DeleteContext: resourceCDROrgDelete,

		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"part_of": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"purge_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCDROrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)

	client, err := config.getFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	org, err := stu3.NewOrganization(config.TimeZone, orgID, name)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	partOf := d.Get("part_of").(string)
	if partOf != "" {
		org.PartOf = &datatypes_go_proto.Reference{
			Reference: &datatypes_go_proto.Reference_OrganizationId{
				OrganizationId: &datatypes_go_proto.ReferenceId{
					Value: partOf,
				},
			},
		}
	}

	// Check if already onboarded
	var onboardedOrg *resources_go_proto.Organization

	onboardedOrg, _, err = client.TenantSTU3.GetOrganizationByID(orgID)
	if err == nil && onboardedOrg != nil {
		d.SetId(onboardedOrg.Id.Value)
		return resourceCDROrgUpdate(ctx, d, m)
	}
	// Do initial boarding
	operation := func() error {
		var resp *cdr.Response
		onboardedOrg, resp, err = client.TenantSTU3.Onboard(org)
		return checkForIAMPermissionErrors(client, resp.Response, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(onboardedOrg.Id.Value)
	return diags
}

func resourceCDROrgRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)

	client, err := config.getFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	org, resp, err := client.TenantSTU3.GetOrganizationByID(orgID)
	if err != nil || resp == nil {
		if resp == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "response is nil",
			})
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

func resourceCDROrgUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)
	id := d.Id()

	client, err := config.getFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	org, _, err := client.TenantSTU3.GetOrganizationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonOrg, err := config.ma.MarshalResource(org)
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
			org.PartOf = &datatypes_go_proto.Reference{
				Reference: &datatypes_go_proto.Reference_OrganizationId{
					OrganizationId: &datatypes_go_proto.ReferenceId{
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

	changedOrg, _ := config.ma.MarshalResource(org)
	patch, err := jsonpatch.DiffBytes(jsonOrg, changedOrg)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = client.OperationsSTU3.Patch("Organization/"+id, patch)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCDROrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	endpoint := d.Get("fhir_store").(string)
	id := d.Id()

	client, err := config.getFHIRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	purgeDelete := d.Get("purge_delete").(bool)

	if !purgeDelete {
		deleted, resp, err := client.OperationsSTU3.Delete(path.Join("Organization", id))
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
	_, resp, err := client.OperationsSTU3.Post(path.Join("$purge"), []byte(``), func(request *http.Request) error {
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
		Refresh:    purgeStateRefreshFunc(client, resp.Header.Get("Location"), id),
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

func purgeStateRefreshFunc(client *cdr.Client, purgeStatusURL, id string) resource.StateRefreshFunc {
	statusURL, err := url.Parse(purgeStatusURL)
	if err != nil {
		return func() (result interface{}, state string, err error) {
			return nil, "FAILED", err
		}
	}
	return func() (interface{}, string, error) {
		contained, resp, err := client.OperationsSTU3.Get(id, func(request *http.Request) error {
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

func checkForIAMPermissionErrors(client iam.TokenRefresher, resp *http.Response, err error) error {
	if resp == nil || resp.StatusCode > 500 {
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		_ = client.TokenRefresh()
		return err
	}
	return backoff.Permanent(err)
}
