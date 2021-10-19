package cdr

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/resources_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceCDRSubscription() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDRSubscriptionCreate,
		ReadContext:   resourceCDRSubscriptionRead,
		UpdateContext: resourceCDRSubscriptionUpdate,
		DeleteContext: resourceCDRSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"fhir_store": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"criteria": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Required: true,
			},
			"headers": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"end": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCDRSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	endpoint := d.Get("endpoint").(string)
	deleteEndpoint := d.Get("delete_endpoint").(string)
	reason := d.Get("reason").(string)
	criteria := d.Get("criteria").(string)
	end := d.Get("end").(string)
	headers := tools.ExpandStringList(d.Get("headers").(*schema.Set).List())
	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	subscription, err := stu3.NewSubscription(
		stu3.WithReason(reason),
		stu3.WithCriteria(criteria),
		stu3.WithHeaders(headers),
		stu3.WithEndpoint(endpoint),
		stu3.WithEndtime(endTime),
		stu3.WithDeleteEndpoint(deleteEndpoint))
	if err != nil {
		return diag.FromErr(err)
	}
	jsonSubscription, err := c.Ma.MarshalResource(subscription)
	if err != nil {
		return diag.FromErr(err)
	}
	var contained *resources_go_proto.ContainedResource

	operation := func() error {
		var resp *cdr.Response
		contained, resp, err = client.OperationsSTU3.Post("Subscription", jsonSubscription)
		return tools.CheckForIAMPermissionErrors(client, resp.Response, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))

	if err != nil {
		return diag.FromErr(fmt.Errorf("create subscription: %w", err))
	}
	createdSub := contained.GetSubscription()
	d.SetId(createdSub.Id.Value)
	return diags
}

func resourceCDRSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription read: %w", err))
	}
	defer client.Close()
	contained, resp, err := client.OperationsSTU3.Get("Subscription/" + d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Errorf("subscription read: %w", err).Error(),
			})
		}
		return diags
	}
	sub := contained.GetSubscription()
	_ = d.Set("endpoint", sub.Channel.Endpoint.Value)
	_ = d.Set("reason", sub.Reason.Value)
	_ = d.Set("criteria", sub.Criteria.Value)
	_ = d.Set("status", sub.Status.Value)
	_ = d.Set("delete_endpoint", stu3.DeleteEndpointValue()(sub))
	headers := make([]string, 0)
	for _, h := range sub.Channel.Header {
		headers = append(headers, h.Value)
	}
	_ = d.Set("headers", headers)
	return diags
}

func resourceCDRSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	id := d.Id()

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	contained, _, err := client.OperationsSTU3.Get("Subscription/" + id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}
	sub := contained.GetSubscription()
	jsonSub, err := c.Ma.MarshalResource(sub)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}
	madeChanges := false

	if d.HasChange("criteria") {
		sub.Criteria.Value = d.Get("criteria").(string)
		madeChanges = true
	}
	if d.HasChange("reason") {
		sub.Reason.Value = d.Get("reason").(string)
		madeChanges = true
	}
	if d.HasChange("endpoint") {
		sub.Channel.Endpoint.Value = d.Get("endpoint").(string)
		madeChanges = true
	}
	if d.HasChange("end") {
		endTime, err := time.Parse(time.RFC3339, d.Get("end").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		sub.End.ValueUs = endTime.UnixNano() / 1000
		madeChanges = true
	}
	if d.HasChange("headers") {
		headers := tools.ExpandStringList(d.Get("headers").(*schema.Set).List())
		sub.Channel.Header = make([]*datatypes_go_proto.String, 0)
		for _, h := range headers {
			sub.Channel.Header = append(sub.Channel.Header, &datatypes_go_proto.String{Value: h})
		}
		madeChanges = true
	}
	if d.HasChange("delete_endpoint") {
		modifyDeleteEndpoint := stu3.WithDeleteEndpoint(d.Get("delete_endpoint").(string))
		if err := modifyDeleteEndpoint(sub); err != nil {
			return diag.FromErr(err)
		}
		madeChanges = true
	}
	if !madeChanges {
		return diags
	}

	changedOrg, _ := c.Ma.MarshalResource(sub)
	patch, err := jsonpatch.DiffBytes(jsonSub, changedOrg)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}
	_, _, err = client.OperationsSTU3.Patch("Subscription/"+id, patch)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}

	return diags
}

func resourceCDRSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	id := d.Id()

	client, err := c.GetFHIRClientFromEndpoint(fhirStore)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	// TODO: Check HTTP 500 issue
	ok, _, err := client.OperationsSTU3.Delete("Subscription/" + id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrDeleteSubscriptionFailed)
	}
	return diags
}
