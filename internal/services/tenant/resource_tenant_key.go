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
				Default:     "8760h", // 365 days
				ForceNew:    true,
				Description: "The duration before the token expires. Uses Go duration format (e.g., '24h', '7d', '1y')",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					duration, err := time.ParseDuration(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid duration, got: %s", key, err))
						return
					}
					if duration <= 0 {
						errs = append(errs, fmt.Errorf("%q must be a positive duration", key))
					}
					return
				},
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

	// Parse the duration string
	duration, err := time.ParseDuration(expirationStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse expiration duration '%s': %w", expirationStr, err)
	}

	apiKey, signature, err := keys.GenerateAPIKey(
		"2",
		signingKey,
		organization,
		environment,
		region,
		project,
		scopes,
		time.Now().Add(duration),
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
	d.SetId(signature)

	return diags
}

// resourceTenantKeyRead reads a tenant key resource
func resourceTenantKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()

	_, signature, err := generateAPIKeyAndSignature(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if signature != id {
		d.SetId("")
		return diags
	}
	return diags
}

// resourceTenantKeyDelete is a stub implementation to resolve the undefined error.
func resourceTenantKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
