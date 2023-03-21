package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceIAMApplication() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    ResourceIAMApplicationV0().CoreConfigSchema().ImpliedType(),
				Upgrade: patchIAMApplicationV0,
				Version: 0,
			},
		},
		SchemaVersion: 1,
		Description:   "Manage HSDP IAM application under a proposition",
		CreateContext: resourceIAMApplicationCreate,
		UpdateContext: resourceIAMApplicationUpdate,
		ReadContext:   resourceIAMApplicationRead,
		DeleteContext: resourceIAMApplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tools.ValidateUpperString,
				ForceNew:     true,
				Description:  "The name of the application.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The description of the application.",
			},
			"proposition_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The proposition ID (GUID) to attach this a application to.",
			},
			"global_reference_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
				Description:      "Reference identifier defined by the provisioning user. Recommend to not set this and let Terraform generate a UUID for you.",
			},
			"wait_for_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Blocks until the application delete has completed. Default: false. The application delete process can take some time as all its associated resources like services and clients are removed recursively. This option is useful for ephemeral environments where the same application might be recreated shortly after a destroy operation.",
			},
		},
	}
}

func resourceIAMApplicationUpdate(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if !d.HasChange("wait_for_delete") {
		return diag.FromErr(fmt.Errorf("only 'wait_for_delete' can be updated"))
	}
	return diags
}

func resourceIAMApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var app iam.Application
	app.Name = d.Get("name").(string) // TODO: this must be all caps
	app.Description = d.Get("description").(string)
	app.PropositionID = d.Get("proposition_id").(string)
	app.GlobalReferenceID = d.Get("global_reference_id").(string)
	if app.GlobalReferenceID == "" {
		result, err := uuid.GenerateUUID()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error generating uuid: %w", err))
		}
		app.GlobalReferenceID = result
	}
	var createdApp *iam.Application
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 8, func() (*http.Response, error) {
		var err error
		createdApp, resp, err = client.Applications.CreateApplication(app)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode() != http.StatusConflict {
			return diag.FromErr(err)
		}
		createdApps, _, err := client.Applications.GetApplications(&iam.GetApplicationsOptions{
			Name:          &app.Name,
			PropositionID: &app.PropositionID,
		})
		if err != nil || len(createdApps) == 0 {
			return diag.FromErr(fmt.Errorf("GetApplications after 409 (len=%d): %w", len(createdApps), err))
		}
		createdApp = createdApps[0]
		if createdApp.Description != app.Description {
			return diag.FromErr(fmt.Errorf("existing application found but description mismatch: '%s' != '%s'", createdApp.Description, app.Description))
		}
		if createdApp.PropositionID != app.PropositionID {
			return diag.FromErr(fmt.Errorf("existing application found but proposition_id mismatch: '%s' != '%s'", createdApp.PropositionID, app.PropositionID))
		}
		// We found a matching existing application, go with it
	}
	if createdApp == nil {
		return diag.FromErr(fmt.Errorf("unexpected failure creating '%s': [%v] [%v]", app.Name, err, resp))
	}
	d.SetId(createdApp.ID)
	return resourceIAMApplicationRead(ctx, d, m)
}

func resourceIAMApplicationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	app, resp, err := client.Applications.GetApplicationByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", app.Name)
	_ = d.Set("description", app.Description)
	_ = d.Set("proposition_id", app.PropositionID)
	_ = d.Set("global_reference_id", app.GlobalReferenceID)
	return diags
}

func resourceIAMApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	app, _, err := client.Applications.GetApplicationByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	waitForDelete := d.Get("wait_for_delete").(bool)

	ok, resp, err := client.Applications.DeleteApplication(*app)
	if err != nil {
		if resp.StatusCode() == http.StatusNotFound { // Gone already
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	if waitForDelete {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"IN_PROGRESS", "QUEUED", "indeterminate"},
			Target:     []string{"SUCCESS"},
			Refresh:    checkAppDeleteStatus(client, id),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: time.Duration(5) * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("waiting for delete: %w", err))
		}
	}
	d.SetId("")
	return diags
}

func checkAppDeleteStatus(client *iam.Client, id string) resource.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		appStatus, resp, err := client.Applications.DeleteStatus(id)
		if err != nil {
			return resp, "FAILED", err
		}
		if appStatus != nil {
			return resp, appStatus.Status, nil
		}
		// We may need to return an error here
		return resp, "IN_PROGRESS", nil
	}
}
