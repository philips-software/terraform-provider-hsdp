package proposition

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var descriptions = map[string]string{
	"proposition": "A proposition represents a deployable solution as a unique identity in a hosting organization. It must have one or more independently manageable applications.",
}

func ResourceIAMProposition() *schema.Resource {
	return &schema.Resource{
		Description: descriptions["proposition"],
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMPropositionCreate,
		ReadContext:   resourceIAMPropositionRead,
		UpdateContext: resourceIAMPropositionUpdate,
		DeleteContext: resourceIAMPropositionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tools.ValidateUpperString,
				ForceNew:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"global_reference_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
			},
			"wait_for_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Blocks until the proposition delete has completed. Default: false. The proposition delete process can take some time as all its associated resources like applications and services are removed recursively. This option is useful for ephemeral environments where the same proposition might be recreated shortly after a destroy operation.",
			},
		},
	}
}

func resourceIAMPropositionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// NOOP
	return resourceIAMPropositionRead(ctx, d, m)
}

func resourceIAMPropositionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var prop iam.Proposition
	prop.Name = d.Get("name").(string) // TODO: this must be all caps
	prop.Description = d.Get("description").(string)
	prop.OrganizationID = d.Get("organization_id").(string)
	prop.GlobalReferenceID = d.Get("global_reference_id").(string)
	if prop.GlobalReferenceID == "" {
		result, err := uuid.GenerateUUID()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error generating uuid: %w", err))
		}
		prop.GlobalReferenceID = result
	}
	var createdProp *iam.Proposition
	var resp *iam.Response

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		createdProp, resp, err = client.Propositions.CreateProposition(prop)
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
		createdProp, _, err = client.Propositions.GetProposition(&iam.GetPropositionsOptions{
			Name:           &prop.Name,
			OrganizationID: &prop.OrganizationID,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("CreateProposition 409 confict, but no match found: %w", err))
		}
		if createdProp.Description != prop.Description {
			return diag.FromErr(fmt.Errorf("existing proposition found but description mismatch: '%s' != '%s'", createdProp.Description, prop.Description))
		}
		if createdProp.OrganizationID != prop.OrganizationID {
			return diag.FromErr(fmt.Errorf("existing proposition found but organization_id mismatch: '%s' != '%s'", createdProp.OrganizationID, prop.OrganizationID))
		}
		// We found a matching existing proposition, go with it
	}
	if createdProp == nil {
		return diag.FromErr(fmt.Errorf("unexpected error creating proposition: %v", resp))
	}
	d.SetId(createdProp.ID)
	return resourceIAMPropositionRead(ctx, d, m)
}

func resourceIAMPropositionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	prop, resp, err := client.Propositions.GetPropositionByID(id)
	if err != nil {
		if errors.Is(err, iam.ErrEmptyResults) || (resp != nil && resp.StatusCode() == http.StatusNotFound) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", prop.Name)
	_ = d.Set("description", prop.Description)
	_ = d.Set("organization_id", prop.OrganizationID)
	_ = d.Set("global_reference_id", prop.GlobalReferenceID)
	return diags
}

func resourceIAMPropositionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	prop, _, err := client.Propositions.GetPropositionByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	waitForDelete := d.Get("wait_for_delete").(bool)

	ok, _, err := client.Propositions.DeleteProposition(*prop)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrInvalidResponse)
	}
	if waitForDelete {
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"IN_PROGRESS", "QUEUED", "indeterminate"},
			Target:     []string{"SUCCESS"},
			Refresh:    checkPropDeleteStatus(client, id),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: time.Duration(5) * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("waiting for delete: %w", err))
		}
	}
	return diags
}

func checkPropDeleteStatus(client *iam.Client, id string) retry.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		propStatus, resp, err := client.Propositions.DeleteStatus(id)
		if err != nil {
			return resp, "FAILED", err
		}
		if propStatus != nil {
			return resp, propStatus.Status, nil
		}
		// We may need to return an error here
		return resp, "IN_PROGRESS", nil
	}
}
