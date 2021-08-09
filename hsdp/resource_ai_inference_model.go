package hsdp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/ai/inference"
)

func resourceAIInferenceModel() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceAIInferenceModelCreate,
		ReadContext:   resourceAIInferenceModelRead,
		DeleteContext: resourceAIInferenceModelDelete,

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
			"version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"artifact_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"environment": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"entry_commands": {
				Type:     schema.TypeList,
				MaxItems: 10,
				MinItems: 1,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeList,
				MaxItems: 20,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compute_environment": {
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
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAIInferenceModelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	version := d.Get("version").(string)
	artifactPath := d.Get("artifact_path").(string)
	entryCommands, _ := collectList("entry_commands", d)
	labels, _ := collectList("labels", d)
	computeEnvironment, _ := collectComputeEnvironment(d)
	sourceCode, _ := collectSourceCode(d)
	additionalConfiguration := d.Get("additional_configuration").(string)

	model := inference.Model{
		ResourceType:            "Model",
		Name:                    name,
		Version:                 version,
		Description:             description,
		ArtifactPath:            artifactPath,
		EntryCommands:           entryCommands,
		ComputeEnvironment:      computeEnvironment,
		SourceCode:              sourceCode,
		AdditionalConfiguration: additionalConfiguration,
		Labels:                  labels,
		Type:                    "sagemaker",
	}
	if v, ok := d.GetOk("environment"); ok {
		vv := v.(map[string]interface{})
		for k, v := range vv {
			model.EnvVars = append(model.EnvVars, inference.EnvironmentVariable{
				Name:  k,
				Value: fmt.Sprint(v),
			})
		}
	}

	var createdModel *inference.Model
	var resp *inference.Response
	// Do initial boarding
	operation := func() error {
		createdModel, resp, err = client.Model.CreateModel(model)
		if resp == nil {
			resp = &inference.Response{}
		}
		return checkForIAMPermissionErrors(client, resp.Response, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdModel.ID)
	return resourceAIInferenceModelRead(ctx, d, m)
}

func collectComputeEnvironment(d *schema.ResourceData) (inference.ReferenceComputeEnvironment, diag.Diagnostics) {
	var diags diag.Diagnostics
	var rce inference.ReferenceComputeEnvironment
	if v, ok := d.GetOk("compute_environment"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			rce = inference.ReferenceComputeEnvironment{
				Reference:  mVi["reference"].(string),
				Identifier: mVi["identifier"].(string),
			}
		}
	}
	return rce, diags
}

func collectSourceCode(d *schema.ResourceData) (inference.SourceCode, diag.Diagnostics) {
	var diags diag.Diagnostics
	var rce inference.SourceCode
	if v, ok := d.GetOk("source_code"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			rce = inference.SourceCode{
				URL:      mVi["url"].(string),
				Branch:   mVi["branch"].(string),
				CommitID: mVi["commit_id"].(string),
				SSHKey:   mVi["ssh_key"].(string),
			}
		}
	}
	return rce, diags
}

func resourceAIInferenceModelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	model, _, err := client.Model.GetModelByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", model.Name)
	_ = d.Set("description", model.Description)
	_ = d.Set("version", model.Version)
	_ = d.Set("artifact_path", model.ArtifactPath)
	_ = d.Set("created", model.Created)
	_ = d.Set("created_by", model.CreatedBy)

	return diags
}

func resourceAIInferenceModelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	resp, err := client.Model.DeleteModel(inference.Model{
		ID: id,
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode == http.StatusNotFound { // Already deleted
			d.SetId("")
			return diags
		}
	}
	d.SetId("")
	return diags
}
