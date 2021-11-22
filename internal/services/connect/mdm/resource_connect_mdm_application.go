package mdm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceMDMApplication() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceMDMApplicationCreate,
		ReadContext:   resourceMDMApplicationRead,
		UpdateContext: resourceMDMApplicationUpdate,
		DeleteContext: resourceMDMApplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tools.ValidateUpperString,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressCaseDiffs),
				ForceNew: true,
			},
			"proposition_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressDefaultSystemValue),
			},
			"global_reference_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaToApplication(d *schema.ResourceData) mdm.Application {
	name := d.Get("name").(string)
	globalReferenceId := d.Get("global_reference_id").(string)
	propositionId := d.Get("proposition_id").(string)

	resource := mdm.Application{
		Name:              name,
		PropositionID:     mdm.Reference{Reference: propositionId},
		GlobalReferenceID: globalReferenceId,
	}
	return resource
}

func applicationToSchema(resource mdm.Application, d *schema.ResourceData) {
	_ = d.Set("name", resource.Name)
	_ = d.Set("guid", resource.ID)
	_ = d.Set("global_reference_id", resource.GlobalReferenceID)
	_ = d.Set("proposition_id", resource.PropositionID.Reference)
}

func resourceMDMApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToApplication(d)

	if resource.GlobalReferenceID == "" {
		result, err := uuid.GenerateUUID()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error generating uuid: %w", err))
		}
		resource.GlobalReferenceID = result
	}
	var created *mdm.Application
	var resp *mdm.Response

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Applications.CreateApplication(resource)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	}, http.StatusForbidden, http.StatusInternalServerError, http.StatusTooManyRequests)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if !(resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusUnprocessableEntity) {
			return diag.FromErr(fmt.Errorf("error creating Application: %w", err))
		}
		found, _, err := client.Applications.GetApplications(&mdm.GetApplicationsOptions{
			Name:          &resource.Name,
			PropositionID: &resource.PropositionID.Reference,
		})
		if err != nil || len(*found) == 0 {
			return diag.FromErr(fmt.Errorf("CreateApplication 409/422 confict, but no match found: %w", err))
		}
		created = &(*found)[0]
		if created.PropositionID.Reference != resource.PropositionID.Reference {
			return diag.FromErr(fmt.Errorf("existing Application found but proposition_id mismatch: '%s' != '%s'", created.PropositionID.Reference, resource.PropositionID.Reference))
		}
		// We found a matching existing Application, go with it
	}
	if created == nil {
		return diag.FromErr(fmt.Errorf("unexpected error creating Application: %v", resp))
	}
	d.SetId(fmt.Sprintf("Application/%s", created.ID))
	_ = d.Set("guid", created.ID)
	return resourceMDMApplicationRead(ctx, d, m)
}

func resourceMDMApplicationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "Application/%s", &id)

	prop, resp, err := client.Applications.GetApplicationByID(id)
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	applicationToSchema(*prop, d)
	return diags
}

func resourceMDMApplicationDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// As Connect MDM does not support MDM Application deletion we simply
	// clear the Application from state. This will be properly implemented
	// once the MDM API balances out
	var diags diag.Diagnostics
	d.SetId("")

	return diags
}

func resourceMDMApplicationUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	resource := schemaToApplication(d)
	updated, _, err := client.Applications.UpdateApplication(resource)
	if err != nil {
		return diag.FromErr(err)
	}
	applicationToSchema(*updated, d)
	return diags
}