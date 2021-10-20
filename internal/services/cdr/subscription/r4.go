package subscription

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	r4dt "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func r4Create(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

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

	subscription, err := r4.NewSubscription(
		r4.WithReason(reason),
		r4.WithCriteria(criteria),
		r4.WithHeaders(headers),
		r4.WithEndpoint(endpoint),
		r4.WithEndtime(endTime),
		r4.WithDeleteEndpoint(deleteEndpoint))
	if err != nil {
		return diag.FromErr(err)
	}
	jsonSubscription, err := c.R4MA.MarshalResource(subscription)
	if err != nil {
		return diag.FromErr(err)
	}
	var contained *r4pb.ContainedResource

	operation := func() error {
		var resp *cdr.Response
		contained, resp, err = client.OperationsR4.Post("Subscription", jsonSubscription)
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

func r4Read(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	contained, resp, err := client.OperationsR4.Get("Subscription/" + d.Id())
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
	_ = d.Set("delete_endpoint", r4.DeleteEndpointValue()(sub))
	headers := make([]string, 0)
	for _, h := range sub.Channel.Header {
		headers = append(headers, h.Value)
	}
	_ = d.Set("headers", headers)
	return diags
}

func r4Update(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()
	contained, _, err := client.OperationsR4.Get("Subscription/" + id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}
	sub := contained.GetSubscription()
	jsonSub, err := c.R4MA.MarshalResource(sub)
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
		sub.Channel.Header = make([]*r4dt.String, 0)
		for _, h := range headers {
			sub.Channel.Header = append(sub.Channel.Header, &r4dt.String{Value: h})
		}
		madeChanges = true
	}
	if d.HasChange("delete_endpoint") {
		modifyDeleteEndpoint := r4.WithDeleteEndpoint(d.Get("delete_endpoint").(string))
		if err := modifyDeleteEndpoint(sub); err != nil {
			return diag.FromErr(err)
		}
		madeChanges = true
	}
	if !madeChanges {
		return diags
	}

	changedOrg, _ := c.R4MA.MarshalResource(sub)
	patch, err := jsonpatch.DiffBytes(jsonSub, changedOrg)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}
	_, _, err = client.OperationsR4.Patch("Subscription/"+id, patch)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}

	return diags
}

func r4Delete(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Check HTTP 500 issue
	id := d.Id()
	ok, _, err := client.OperationsR4.Delete("Subscription/" + id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrDeleteSubscriptionFailed)
	}
	return diags
}
