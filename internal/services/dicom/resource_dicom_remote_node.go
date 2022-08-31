package dicom

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/dicom"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceDICOMRemoteNode() *schema.Resource {
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
				Required: true,
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
				Elem:     networkConnectionSchema(),
			},
		},
	}
}

func networkConnectionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_secure": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disable_ipv6": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3000,
				ForceNew: true,
			},
			// ---Advanced features start
			"pdu_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  65535,
				ForceNew: true,
			},
			"artim_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3000,
				ForceNew: true,
			},
			"association_idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4500,
				ForceNew: true,
			},
		},
	}
}

func resourceDICOMRemoteNodeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	organizationID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	operation := func() error {
		var resp *dicom.Response
		_, resp, err = client.Config.DeleteRemoteNode(dicom.RemoteNode{ID: d.Id()}, &dicom.QueryOptions{
			OrganizationID: &organizationID,
		})
		return tools.CheckForPermissionErrors(client, resp, err)
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
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	organizationID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	var node *dicom.RemoteNode
	operation := func() error {
		var resp *dicom.Response
		node, resp, err = client.Config.GetRemoteNode(d.Id(), &dicom.QueryOptions{
			OrganizationID: &organizationID,
		})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))

	if err != nil {
		if errors.Is(err, dicom.ErrNotFound) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("title", node.Title)
	_ = d.Set("ae_title", node.AETitle)
	return diags
}

func resourceDICOMRemoteNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	configURL := d.Get("config_url").(string)
	organizationID := d.Get("organization_id").(string)
	client, err := c.GetDICOMConfigClient(configURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	node := dicom.RemoteNode{
		Title:   d.Get("title").(string),
		AETitle: d.Get("ae_title").(string),
	}
	if v, ok := d.GetOk("network_connection"); ok {
		vL := v.(*schema.Set).List()
		networkConnection := dicom.NetworkConnection{}
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			networkConnection.HostName = mVi["hostname"].(string)
			networkConnection.IPAddress = mVi["ip_address"].(string)
			networkConnection.DisableIPv6 = mVi["disable_ipv6"].(bool)

			networkConnection.Port = mVi["port"].(int)
			networkConnection.NetworkTimeout = mVi["network_timeout"].(int)
			networkConnection.AdvancedSettings = &dicom.AdvancedSettings{
				PDULength:              mVi["pdu_length"].(int),
				ArtimTimeout:           mVi["artim_timeout"].(int),
				AssociationIdleTimeout: mVi["association_idle_timeout"].(int),
			}
		}
		node.NetworkConnection = networkConnection
	}

	var created *dicom.RemoteNode
	operation := func() error {
		var resp *dicom.Response
		created, resp, err = client.Config.CreateRemoteNode(node, &dicom.QueryOptions{
			OrganizationID: &organizationID,
		})
		return tools.CheckForPermissionErrors(client, resp, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}
	if created == nil || created.ID == "" {
		return diag.FromErr(fmt.Errorf("failed to create remote node, even though no error was reported"))
	}
	d.SetId(created.ID)
	return resourceDICOMRemoteNodeRead(ctx, d, m)
}
