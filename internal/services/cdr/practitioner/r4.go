package practitioner

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cenkalti/backoff"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	jsonpatch "github.com/herkyl/patchwerk"
	"github.com/philips-software/go-hsdp-api/cdr"
	pr "github.com/philips-software/go-hsdp-api/cdr/helper/fhir/r4/practitioner"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

type identifier struct {
	System string
	Value  string
	Use    string
}

type name struct {
	Text   string
	Given  []string
	Family string
}

func schemaToName(d *schema.ResourceData) []name {
	var resources []name
	if v, ok := d.GetOk("name"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			resources = append(resources, name{
				Text:   mVi["text"].(string),
				Family: mVi["family"].(string),
				Given:  tools.ExpandStringList(mVi["given"].(*schema.Set).List()),
			})
		}
	}
	return resources
}

func schemaToIdentifier(d *schema.ResourceData) []identifier {
	var resources []identifier
	if v, ok := d.GetOk("identifier"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			resources = append(resources, identifier{
				System: mVi["system"].(string),
				Value:  mVi["value"].(string),
				Use:    mVi["use"].(string),
			})
		}
	}
	return resources
}

func r4Create(ctx context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names := schemaToName(d)
	identifiers := schemaToIdentifier(d)

	resource, err := pr.NewPractitioner()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, i := range identifiers {
		if pr.WithIdentifier(i.System, i.Value, i.Use)(resource) != nil {
			return diag.FromErr(err)
		}
	}
	for _, n := range names {
		if pr.WithName(n.Text, n.Family, n.Given)(resource) != nil {
			return diag.FromErr(err)
		}
	}
	jsonResource, err := c.R4MA.MarshalResource(resource)
	if err != nil {
		return diag.FromErr(err)
	}
	var contained *r4pb.ContainedResource

	operation := func() error {
		var resp *cdr.Response
		contained, resp, err = client.OperationsR4.Post("Practitioner", jsonResource)
		if resp == nil {
			if err != nil {
				return err
			}
			return fmt.Errorf("OperationsR4.Post: response is nil")
		}
		return tools.CheckForIAMPermissionErrors(client, resp.Response, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))

	if err != nil {
		return diag.FromErr(fmt.Errorf("create practitioner: %w", err))
	}
	createdResource := contained.GetPractitioner()
	d.SetId(createdResource.Id.Value)
	return r4Read(ctx, c, client, d, m)
}

func r4Read(_ context.Context, _ *config.Config, client *cdr.Client, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	contained, resp, err := client.OperationsR4.Get("Practitioner/" + d.Id())
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone) {
			d.SetId("")
			return diags
		}
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Errorf("practitioner read: %w", err).Error(),
			})
		}
		return diags
	}
	resource := contained.GetPractitioner()

	// Set identifier
	a := &schema.Set{F: schema.HashResource(identifierSchema())}
	for _, cc := range resource.Identifier {
		entry := make(map[string]interface{})
		if cc.System != nil {
			entry["system"] = cc.System.String()
		}
		if cc.Value != nil {
			entry["value"] = cc.Value.String()
		}
		if cc.Use != nil {
			entry["use"] = strings.ToLower(cc.Use.String())
		}
		a.Add(entry)
	}

	// Set names
	n := &schema.Set{F: schema.HashResource(nameSchema())}
	for _, cc := range resource.Name {
		entry := make(map[string]interface{})
		var gg []string
		for _, g := range cc.Given {
			gg = append(gg, g.String())
		}
		if cc.Family != nil {
			entry["family"] = cc.Family.String()
		}
		if cc.Text != nil {
			entry["text"] = cc.Text.String()
		}
		entry["given"] = tools.SchemaSetStrings(gg)
		n.Add(entry)
	}

	// Set meta
	if resource.Meta != nil {
		if resource.Meta.VersionId != nil {
			_ = d.Set("version_id", resource.Meta.VersionId.String())
		}
		if resource.Meta.LastUpdated != nil {
			_ = d.Set("last_updated", resource.Meta.LastUpdated.String())
		}
	}
	return diags
}

func r4Update(_ context.Context, c *config.Config, client *cdr.Client, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()
	contained, _, err := client.OperationsR4.Get("Practitioner/" + id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("practitioner update: %w", err))
	}
	resource := contained.GetPractitioner()
	jsonResource, err := c.R4MA.MarshalResource(resource)
	if err != nil {
		return diag.FromErr(fmt.Errorf("practitioner update: %w", err))
	}
	madeChanges := false

	if d.HasChange("identifier") {
		identifiers := schemaToIdentifier(d)
		resource.Identifier = nil
		for _, i := range identifiers {
			if pr.WithIdentifier(i.System, i.Value, i.Use)(resource) != nil {
				return diag.FromErr(err)
			}
		}
		madeChanges = true
	}
	if d.HasChange("name") {
		names := schemaToName(d)
		for _, n := range names {
			if pr.WithName(n.Text, n.Family, n.Given)(resource) != nil {
				return diag.FromErr(err)
			}
		}
		madeChanges = true
	}

	if !madeChanges {
		return diags
	}
	resource.Meta = nil

	changedResource, _ := c.R4MA.MarshalResource(resource)
	patch, err := jsonpatch.DiffBytes(jsonResource, changedResource)
	if err != nil {
		return diag.FromErr(fmt.Errorf("practitioner update: %w", err))
	}
	_, _, err = client.OperationsR4.Patch("Practitioner/"+id, patch)
	if err != nil {
		return diag.FromErr(fmt.Errorf("practitioner update: %w", err))
	}

	return diags
}

func r4Delete(_ context.Context, _ *config.Config, client *cdr.Client, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Check HTTP 500 issue
	id := d.Id()
	ok, _, err := client.OperationsR4.Delete("Practitioner/" + id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(config.ErrDeleteSubscriptionFailed)
	}
	return diags
}
