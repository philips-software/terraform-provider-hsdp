package mdm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMDMProposition() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		CreateContext: resourceMDMPropositionCreate,
		ReadContext:   resourceMDMPropositionRead,
		UpdateContext: resourceMDMPropositionUpdate,
		DeleteContext: resourceMDMPropositionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tools.ValidateUpperString,
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressCaseDiffs),
				ForceNew: true,
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
				DiffSuppressFunc: tools.SuppressMulti(
					tools.SuppressWhenGenerated,
					tools.SuppressDefaultSystemValue),
			},
			"global_reference_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tools.SuppressWhenGenerated,
			},
			"status": {
				Type:         schema.TypeString,
				ValidateFunc: tools.ValidateUpperString,
				Required:     true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proposition_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaToProposition(d *schema.ResourceData) mdm.Proposition {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	status := d.Get("status").(string)
	organizationId := d.Get("organization_id").(string)

	resource := mdm.Proposition{
		Name:        name,
		Status:      status,
		Description: description,
	}
	if len(organizationId) > 0 {
		identifier := mdm.Identifier{}
		parts := strings.Split(organizationId, "|")
		if len(parts) > 1 {
			identifier.System = parts[0]
			identifier.Value = parts[1]
		} else {
			identifier.Value = organizationId
		}
		resource.OrganizationGuid = identifier
	}
	return resource
}

func propositionToSchema(resource mdm.Proposition, d *schema.ResourceData) {
	_ = d.Set("description", resource.Description)
	_ = d.Set("name", resource.Name)
	_ = d.Set("guid", resource.ID)
	_ = d.Set("status", resource.Status)
	_ = d.Set("proposition_guid", resource.PropositionGuid)

	if resource.OrganizationGuid.Value != "" {
		value := resource.OrganizationGuid.Value
		if resource.OrganizationGuid.System != "" {
			value = fmt.Sprintf("%s|%s", resource.OrganizationGuid.System, resource.OrganizationGuid.Value)
		}
		_ = d.Set("organization_id", value)
	}
}

func resourceMDMPropositionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := schemaToProposition(d)

	if resource.GlobalReferenceID == "" {
		result, err := uuid.GenerateUUID()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error generating uuid: %w", err))
		}
		resource.GlobalReferenceID = result
	}
	var created *mdm.Proposition
	var resp *mdm.Response

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		created, resp, err = client.Propositions.CreateProposition(resource)
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
		if !(resp.StatusCode() == http.StatusConflict || resp.StatusCode() == http.StatusUnprocessableEntity) {
			return diag.FromErr(fmt.Errorf("error creating Proposition (%d): %w", resp.StatusCode(), err))
		}
		found, _, getErr := client.Propositions.GetProposition(&mdm.GetPropositionsOptions{
			Name:           &resource.Name,
			OrganizationID: &resource.OrganizationGuid.Value,
		})
		if getErr != nil {
			return diag.FromErr(fmt.Errorf("CreateProposition 409/422 confict, but no match found: %w", err))
		}
		if found.Description != resource.Description {
			return diag.FromErr(fmt.Errorf("existing proposition found but description mismatch: '%s' != '%s'", found.Description, resource.Description))
		}
		if found.OrganizationGuid.Value != resource.OrganizationGuid.Value {
			return diag.FromErr(fmt.Errorf("existing proposition found but organization_id mismatch: '%s' != '%s'", found.OrganizationGuid.Value, resource.OrganizationGuid.Value))
		}
		// We found a matching existing proposition, go with it
		created = found
	}
	if created == nil {
		return diag.FromErr(fmt.Errorf("unexpected error creating proposition: %v", resp))
	}
	d.SetId(fmt.Sprintf("Proposition/%s", created.ID))
	propositionToSchema(*created, d)
	return diags
}

func resourceMDMPropositionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var id string
	_, _ = fmt.Sscanf(d.Id(), "Proposition/%s", &id)
	var resource *mdm.Proposition
	var resp *mdm.Response
	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		resource, resp, err = client.Propositions.GetPropositionByID(id)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		if errors.Is(err, mdm.ErrEmptyResult) || (resp != nil && (resp.StatusCode() == http.StatusNotFound || resp.StatusCode() == http.StatusGone)) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	propositionToSchema(*resource, d)
	return diags
}

func resourceMDMPropositionDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// As Connect MDM does not support MDM proposition deletion we simply
	// clear the proposition from state. This will be properly implemented
	// once the MDM API balances out
	var diags diag.Diagnostics
	d.SetId("")

	return diags
}

func resourceMDMPropositionUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	client, err := c.MDMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Get("guid").(string)

	if !d.HasChange("status") {
		return diag.FromErr(fmt.Errorf("only 'status' can be updated, this is a provider bug"))
	}

	resource, _, err := client.Propositions.GetPropositionByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	resource.Status = d.Get("status").(string)
	updated, _, err := client.Propositions.UpdateProposition(*resource)
	if err != nil {
		return diag.FromErr(err)
	}
	propositionToSchema(*updated, d)
	return diags
}
