package inference

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceAIInferenceComputeEnvironment() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceAIInferenceComputeEnvironmentCreate,
		ReadContext:   resourceAIInferenceComputeEnvironmentRead,
		DeleteContext: resourceAIInferenceComputeEnvironmentDelete,

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
			"image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_factory": {
				Type:     schema.TypeBool,
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

func resourceAIInferenceComputeEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	image := d.Get("image").(string)

	var createdEnv *ai.ComputeEnvironment
	// Do initial boarding
	operation := func() error {
		var resp *ai.Response
		createdEnv, resp, err = client.ComputeEnvironment.CreateComputeEnvironment(ai.ComputeEnvironment{
			ResourceType: "ComputeEnvironment",
			Name:         name,
			Description:  description,
			Image:        image,
		})
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
	d.SetId(createdEnv.ID)
	return resourceAIInferenceComputeEnvironmentRead(ctx, d, m)
}

func resourceAIInferenceComputeEnvironmentRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var resp *ai.Response
	var diags diag.Diagnostics
	var err error

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(fmt.Errorf("AIInferenceComputeEnvironmentRead: %w", err))
	}
	defer client.Close()

	id := d.Id()

	var env *ai.ComputeEnvironment
	err = tools.TryAICall(func() (*ai.Response, error) {
		var err error
		env, resp, err = client.ComputeEnvironment.GetComputeEnvironmentByID(id)
		return resp, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", env.Name)
	_ = d.Set("description", env.Description)
	_ = d.Set("image", env.Image)
	_ = d.Set("is_factory", env.IsFactory)
	_ = d.Set("created", env.Created)
	_ = d.Set("created_by", env.CreatedBy)
	_ = d.Set("reference", fmt.Sprintf("%s/%s", env.ResourceType, env.ID))

	return diags
}

func resourceAIInferenceComputeEnvironmentDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	resp, err := client.ComputeEnvironment.DeleteComputeEnvironment(ai.ComputeEnvironment{
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
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
