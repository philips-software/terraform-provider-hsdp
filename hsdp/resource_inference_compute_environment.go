package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/inference"
)

func resourceInferenceComputeEnvironment() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceInferenceComputeEnvironmentCreate,
		ReadContext:   resourceInferenceComputeEnvironmentRead,
		DeleteContext: resourceInferenceComputeEnvironmentDelete,

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
		},
	}
}

func resourceInferenceComputeEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	image := d.Get("image").(string)

	createdEnv, _, err := client.ComputeEnvironment.CreateComputeEnvironment(inference.ComputeEnvironment{
		Name:        name,
		Description: description,
		Image:       image,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdEnv.ID)
	return resourceInferenceComputeEnvironmentRead(ctx, d, m)
}

func resourceInferenceComputeEnvironmentRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	env, _, err := client.ComputeEnvironment.GetComputeEnvironmentByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", env.Name)
	_ = d.Set("description", env.Description)
	_ = d.Set("image", env.Image)
	_ = d.Set("is_factory", env.IsFactory)
	_ = d.Set("created", env.Created)
	_ = d.Set("created_by", env.CreatedBy)

	return diags
}

func resourceInferenceComputeEnvironmentDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	_, err = client.ComputeEnvironment.DeleteComputeEnvironment(inference.ComputeEnvironment{
		ID: id,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
