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

func ResourceAIInferenceComputeTarget() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceAIInferenceComputeTargetCreate,
		ReadContext:   resourceAIInferenceComputeTargetRead,
		DeleteContext: resourceAIInferenceComputeTargetDelete,

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
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage": {
				Type:     schema.TypeInt,
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

func resourceAIInferenceComputeTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	instanceType := d.Get("instance_type").(string)
	storage := d.Get("storage").(int)

	var createdTarget *ai.ComputeTarget
	var resp *ai.Response
	// Do initial boarding
	operation := func() error {
		createdTarget, resp, err = client.ComputeTarget.CreateComputeTarget(ai.ComputeTarget{
			ResourceType: "ComputeTarget",
			Name:         name,
			Description:  description,
			InstanceType: instanceType,
			Storage:      storage,
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
	d.SetId(createdTarget.ID)
	return resourceAIInferenceComputeTargetRead(ctx, d, m)
}

func resourceAIInferenceComputeTargetRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	target, _, err := client.ComputeTarget.GetComputeTargetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", target.Name)
	_ = d.Set("description", target.Description)
	_ = d.Set("instance_type", target.InstanceType)
	_ = d.Set("storage", target.Storage)
	_ = d.Set("is_factory", target.IsFactory)
	_ = d.Set("created", target.Created)
	_ = d.Set("created_by", target.CreatedBy)
	_ = d.Set("reference", fmt.Sprintf("%s/%s", target.ResourceType, target.ID))

	return diags
}

func resourceAIInferenceComputeTargetDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)
	endpoint := d.Get("endpoint").(string)
	client, err := c.GetAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	resp, err := client.ComputeTarget.DeleteComputeTarget(ai.ComputeTarget{
		ID: id,
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode() == http.StatusNotFound { // Already deleted
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
