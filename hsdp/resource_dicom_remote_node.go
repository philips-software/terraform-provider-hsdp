package hsdp

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
)

func resourceDICOMRemoteNode() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceDICOMRemoteNodeCreate,
		ReadContext:   resourceDICOMRemoteNodeRead,
		DeleteContext: resourceDICOMRemoteNodeDelete,

		Schema: map[string]*schema.Schema{
			"config_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": { // Query
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ae_title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_connection": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_secure": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disable_ipv6": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  104,
						},
						// ---Advanced features start
						"pdu_length": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"artim_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"association_idle_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"network_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceDICOMRemoteNodeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	operation := func() error {
		var resp *dicom.Response
		_, resp, err = client.Config.DeleteRemoteNode(dicom.RemoteNode{ID: d.Id()})
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceDICOMRemoteNodeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	var node *dicom.RemoteNode
	operation := func() error {
		var resp *dicom.Response
		node, resp, err = client.Config.GetRemoteNode(d.Id(), nil)
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("title", node.Title)
	_ = d.Set("ae_title", node.AETitle)
	// TODO: set other field
	return diags
}

func resourceDICOMRemoteNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	configURL := d.Get("config_url").(string)
	client, err := config.getDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	node := dicom.RemoteNode{
		Title:   d.Get("title").(string),
		AETitle: d.Get("ae_title").(string),
	}
	var created *dicom.RemoteNode
	operation := func() error {
		var resp *dicom.Response
		created, resp, err = client.Config.CreateRemoteNode(node)
		return checkForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(created.ID)
	return resourceDICOMRemoteNodeRead(ctx, d, m)
}
