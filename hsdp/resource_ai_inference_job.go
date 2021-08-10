package hsdp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/ai"
)

func resourceAIInferenceJob() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceAIInferenceJobCreate,
		ReadContext:   resourceAIInferenceJobRead,
		DeleteContext: resourceAIInferenceJobDelete,

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
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
				ForceNew: true,
			},
			"environment": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"command_args": {
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
			"model": {
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
			"input": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 0,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"output": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 0,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"completed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
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

func resourceAIInferenceJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	commandArgs, _ := collectList("command_args", d)
	labels, _ := collectList("labels", d)
	computeTarget, _ := collectComputeTarget(d)
	computeModel, _ := collectComputeModel(d)
	inputs, _ := collectInputs(d)
	outputs, _ := collectOutputs(d)
	timeout := d.Get("timeout").(int)

	job := ai.Job{
		ResourceType:  "InferenceJob",
		Name:          name,
		Description:   description,
		CommandArgs:   commandArgs,
		ComputeTarget: computeTarget,
		Model:         computeModel,
		Input:         inputs,
		Output:        outputs,
		Timeout:       timeout,
		Labels:        labels,
		Type:          "sagemaker",
	}
	if v, ok := d.GetOk("environment"); ok {
		vv := v.(map[string]interface{})
		for k, v := range vv {
			job.EnvVars = append(job.EnvVars, ai.EnvironmentVariable{
				Name:  k,
				Value: fmt.Sprint(v),
			})
		}
	}

	var createdJob *ai.Job
	var resp *ai.Response
	// Do initial boarding
	operation := func() error {
		createdJob, resp, err = client.Job.CreateJob(job)
		if resp == nil {
			resp = &ai.Response{}
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
	d.SetId(createdJob.ID)
	return resourceAIInferenceJobRead(ctx, d, m)
}

func collectComputeTarget(d *schema.ResourceData) (ai.ReferenceComputeTarget, diag.Diagnostics) {
	var diags diag.Diagnostics
	var target ai.ReferenceComputeTarget
	if v, ok := d.GetOk("compute_target"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			target = ai.ReferenceComputeTarget{
				Reference:  mVi["reference"].(string),
				Identifier: mVi["identifier"].(string),
			}
		}
	}
	return target, diags
}

func collectComputeModel(d *schema.ResourceData) (ai.ReferenceComputeModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model ai.ReferenceComputeModel
	if v, ok := d.GetOk("model"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			model = ai.ReferenceComputeModel{
				Reference:  mVi["reference"].(string),
				Identifier: mVi["identifier"].(string),
			}
		}
	}
	return model, diags
}

func collectInputs(d *schema.ResourceData) ([]ai.InputEntry, diag.Diagnostics) {
	var diags diag.Diagnostics
	var inputs []ai.InputEntry
	if v, ok := d.GetOk("input"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			input := ai.InputEntry{
				URL:  mVi["url"].(string),
				Name: mVi["name"].(string),
			}
			inputs = append(inputs, input)
		}
	}
	return inputs, diags
}

func collectOutputs(d *schema.ResourceData) ([]ai.OutputEntry, diag.Diagnostics) {
	var diags diag.Diagnostics
	var outputs []ai.OutputEntry
	if v, ok := d.GetOk("output"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			input := ai.OutputEntry{
				URL:  mVi["url"].(string),
				Name: mVi["name"].(string),
			}
			outputs = append(outputs, input)
		}
	}
	return outputs, diags
}

func resourceAIInferenceJobRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	job, _, err := client.Job.GetJobByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", job.Name)
	_ = d.Set("description", job.Description)
	_ = d.Set("timeout", job.Timeout)
	_ = d.Set("duration", job.Duration)
	_ = d.Set("completed", job.ComputeTarget)
	_ = d.Set("status", job.Status)
	_ = d.Set("command_args", job.CommandArgs)
	_ = d.Set("created", job.Created)
	_ = d.Set("created_by", job.CreatedBy)
	_ = d.Set("reference", fmt.Sprintf("%s/%s", job.ResourceType, job.ID))

	return diags
}

func resourceAIInferenceJobDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	job := ai.Job{
		ID: id,
	}

	_, _ = client.Job.TerminateJob(job) // Just to be sure

	resp, err := client.Job.DeleteJob(job)
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
