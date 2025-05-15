package tenant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/loafoe/caddy-token/keys"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceTenantKey() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceTenantKeyCreate,
		ReadContext:   resourceTenantKeyRead,
		DeleteContext: resourceTenantKeyDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"signing_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"expiration": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC3339), // 1 year from now
				ForceNew:    true,
				Description: "Expiration time in RFC3339 format (e.g., '2025-12-31T23:59:59Z')",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := time.Parse(time.RFC3339, v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid RFC3339 time string (e.g., '2025-12-31T23:59:59Z'), got: %s", key, err))
						return
					}
					return
				},
			},
			"salt": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "Salt value to use for generating a deterministic API key",
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "us-east",
				ForceNew: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "prod",
				ForceNew: true,
			},
			"key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"signature": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The signature of the generated API key",
			},
		},
	}
}

// generateAPIKeyAndSignature is a helper function to generate an API key and extract its signature
func generateAPIKeyAndSignature(d *schema.ResourceData) (string, string, error) {
	project := d.Get("project").(string)
	organization := d.Get("organization").(string)
	signingKey := d.Get("signing_key").(string)
	scopes := tools.ExpandStringList(d.Get("scopes").(*schema.Set).List())
	region := d.Get("region").(string)
	environment := d.Get("environment").(string)
	expirationStr := d.Get("expiration").(string)
	salt := d.Get("salt").(string)

	// Parse the expiration time string
	expirationTime, err := time.Parse(time.RFC3339, expirationStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse expiration time '%s': %w", expirationStr, err)
	}

	// Use GenerateDeterministicAPIKey for a deterministic, non-time-dependent key
	apiKey, signature, err := keys.GenerateDeterministicAPIKey(
		"2",
		signingKey,
		organization,
		environment,
		region,
		project,
		scopes,
		expirationTime,
		salt,
	)
	if err != nil {
		return "", "", err
	}
	return apiKey, signature, nil
}

// resourceTenantKeyCreate creates a tenant key resource
func resourceTenantKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	apiKey, signature, err := generateAPIKeyAndSignature(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("key", apiKey)
	d.Set("signature", signature)
	d.SetId(signature)

	return diags
}

// resourceTenantKeyRead reads a tenant key resource
func resourceTenantKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()

	apiKey, signature, err := generateAPIKeyAndSignature(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if signature != id {
		d.SetId("")
		return diags
	}

	d.Set("key", apiKey)
	d.Set("signature", signature)
	return diags
}

// resourceTenantKeyDelete is a stub implementation to resolve the undefined error.
func resourceTenantKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
