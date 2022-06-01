package subscription

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

func stu3Create(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	jsonSubscription, err := c.STU3MA.MarshalResource(subscription)
	if err != nil {
		return diag.FromErr(err)
	}
	var contained *resources_go_proto.ContainedResource

	operation := func() error {
		var resp *cdr.Response
		contained, resp, err = client.OperationsSTU3.Post("Subscription", jsonSubscription)
		if resp == nil {
			if err != nil {
				return err
			}
			return fmt.Errorf("OperationsSTU3.Post: response is nil")
		}
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

func stu3Read(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	contained, resp, err := client.OperationsSTU3.Get("Subscription/" + d.Id())
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
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

func stu3Update(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()
	contained, _, err := client.OperationsSTU3.Get("Subscription/" + id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subscription update: %w", err))
	}
	sub := contained.GetSubscription()
	jsonSub, err := c.STU3MA.MarshalResource(sub)
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

	changedOrg, _ := c.STU3MA.MarshalResource(sub)
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

func stu3Delete(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Check HTTP 500 issue
	id := d.Id()
	ok, _, err := client.OperationsSTU3.Delete("Subscription/" + id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrDeleteSubscriptionFailed)
	}
	return diags
}
