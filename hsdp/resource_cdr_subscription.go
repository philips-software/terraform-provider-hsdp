package hsdp

import (
	"context"
	"github.com/google/fhir/go/proto/google/fhir/proto/stu3/datatypes_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr/helper/fhir/stu3"
	"time"
)

func resourceCDRSubscription() *schema.Resource {
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
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"criteria": {
				Type:     schema.TypeString,
				Required: true,
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
	config := m.(*Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)
	endpoint := d.Get("endpoint").(string)
	reason := d.Get("reason").(string)
	criteria := d.Get("criteria").(string)
	end := d.Get("end").(string)
	headers := expandStringList(d.Get("headers").(*schema.Set).List())
	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := config.getFHIRClient(fhirStore, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	subscription, err := stu3.NewSubscription(
		stu3.WithReason(reason),
		stu3.WithCriteria(criteria),
		stu3.WithHeaders(headers),
		stu3.WithEndpoint(endpoint),
		stu3.WithEndtime(endTime))
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	jsonSubscription, err := config.ma.MarshalResource(subscription)
	if err != nil {
		return diag.FromErr(err)
	}
	contained, resp, err := client.OperationsSTU3.Post("Subscription", jsonSubscription)
	if err != nil || resp == nil {
		return diag.FromErr(err)
	}
	createdSub := contained.GetSubscription()
	d.SetId(createdSub.Id.Value)
	return diags
}

func resourceCDRSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)

	client, err := config.getFHIRClient(fhirStore, orgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	contained, resp, err := client.OperationsSTU3.Get("Subscription/" + d.Id())
	if err != nil || resp == nil {
		if resp == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "response is nil",
			})
		}
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
		return diags
	}
	sub := contained.GetSubscription()
	_ = d.Set("endpoint", sub.Channel.Endpoint.Value)
	_ = d.Set("reason", sub.Reason.Value)
	_ = d.Set("criteria", sub.Criteria.Value)
	_ = d.Set("status", sub.Status.Value)
	headers := make([]string, 0)
	for _, h := range sub.Channel.Header {
		headers = append(headers, h.Value)
	}
	_ = d.Set("headers", headers)
	return diags
}

func resourceCDRSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)
	id := d.Id()

	client, err := config.getFHIRClient(fhirStore, orgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	contained, _, err := client.OperationsSTU3.Get("Subscription/" + id)
	if err != nil {
		return diag.FromErr(err)
	}
	sub := contained.GetSubscription()
	jsonSub, err := config.ma.MarshalResource(sub)
	if err != nil {
		return diag.FromErr(err)
	}
	madeChanges := false

	if d.HasChange("criteria") {
		sub.Criteria.Value = d.Get("name").(string)
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
		headers := expandStringList(d.Get("headers").(*schema.Set).List())
		sub.Channel.Header = make([]*datatypes_go_proto.String, 0)
		for _, h := range headers {
			sub.Channel.Header = append(sub.Channel.Header, &datatypes_go_proto.String{Value: h})
		}
		madeChanges = true
	}
	if !madeChanges {
		return diags
	}

	changedOrg, _ := config.ma.MarshalResource(sub)
	patch, err := jsonpatch.DiffBytes(jsonSub, changedOrg)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = client.OperationsSTU3.Patch("Subscription/"+id, patch)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCDRSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	fhirStore := d.Get("fhir_store").(string)
	orgID := d.Get("org_id").(string)
	id := d.Id()

	client, err := config.getFHIRClient(fhirStore, orgID)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	ok, _, err := client.OperationsSTU3.Delete("Subscription/" + id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(ErrDeleteSubscriptionFailed)
	}
	return diags
}
