package ai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/go-hsdp-api/ai/workspace"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceAIWorkspace() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceAIWorkspaceCreate,
		ReadContext:   resourceAIWorkspaceRead,
		DeleteContext: resourceAIWorkspaceDelete,

		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeList,
				MaxItems: 20,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compute_target": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				MinItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reference": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"identifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"source_code": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				MinItems: 0,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"commit_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ssh_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"additional_configuration": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAIWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIWorkspaceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	labels, _ := tools.CollectList("labels", d)
	computeTarget, _ := collectComputeTarget(d)
	sourceCode, _ := collectSourceCode(d)
	additionalConfiguration := d.Get("additional_conifguration").(string)

	ws := workspace.Workspace{
		ResourceType:            "Workspace",
		Name:                    name,
		Description:             description,
		ComputeTarget:           computeTarget,
		SourceCode:              sourceCode,
		AdditionalConfiguration: additionalConfiguration,
		Labels:                  labels,
		Type:                    "sagemaker",
	}

	var createdWorkspace *workspace.Workspace
	var resp *ai.Response
	// Do initial boarding
	operation := func() error {
		createdWorkspace, resp, err = client.Workspace.CreateWorkspace(ws)
		if resp == nil {
			resp = &ai.Response{}
		}
		return tools.CheckForIAMPermissionErrors(client, resp.Response, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdWorkspace.ID)
	return resourceAIWorkspaceRead(ctx, d, m)
}

func resourceAIWorkspaceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIWorkspaceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	ws, _, err := client.Workspace.GetWorkspaceByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", ws.Name)
	_ = d.Set("description", ws.Description)
	_ = d.Set("additional_configuration", ws.AdditionalConfiguration)
	_ = d.Set("created", ws.Created)
	_ = d.Set("created_by", ws.CreatedBy)
	_ = d.Set("reference", fmt.Sprintf("%s/%s", ws.ResourceType, ws.ID))
	return diags
}

func resourceAIWorkspaceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIWorkspaceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	ws := workspace.Workspace{
		ID: id,
	}

	resp, err := client.Workspace.DeleteWorkspace(ws)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode == http.StatusNotFound { // Already deleted
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
