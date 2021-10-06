package hsdp

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-software/go-hsdp-api/iam"
)

func resourceIAMService() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 3,
		CreateContext: resourceIAMServiceCreate,
		ReadContext:   resourceIAMServiceRead,
		UpdateContext: resourceIAMServiceUpdate,
		DeleteContext: resourceIAMServiceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressCaseDiffs,
			},
			"description": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"validity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      12,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 600),
			},
			"self_managed_private_key": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				ConflictsWith: []string{"self_managed_certificate"},
			},
			"self_managed_certificate": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"self_managed_private_key", "expires_on"},
			},
			"private_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_on": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressWhenGenerated,
				ConflictsWith:    []string{"self_managed_certificate"},
			},
			"scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 1, // openid
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_scopes": {
				Type:     schema.TypeSet,
				MaxItems: 100,
				MinItems: 1, // openid
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIAMServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var s iam.Service
	s.Description = d.Get("description").(string)
	s.Name = d.Get("name").(string)
	s.ApplicationID = d.Get("application_id").(string)
	s.Validity = d.Get("validity").(int)
	scopes := expandStringList(d.Get("scopes").(*schema.Set).List())
	defaultScopes := expandStringList(d.Get("default_scopes").(*schema.Set).List())
	expiresOn := d.Get("expires_on").(string)
	selfPrivateKey := d.Get("self_managed_private_key").(string)
	selfCertificate := d.Get("self_managed_certificate").(string)
	if selfPrivateKey == "" && expiresOn != "" {
		return diag.FromErr(fmt.Errorf("you cannot set an 'expires_on' value without also specifying the 'self_managed_private_key'"))
	}

	createdService, _, err := client.Services.CreateService(s)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set certificate if set from the get go
	if selfPrivateKey != "" || selfCertificate != "" {
		diags = setSelfManaged(client, *createdService, d)
		if len(diags) > 0 {
			_, _, _ = client.Services.DeleteService(*createdService) // Cleanup
			return diags
		}
		_ = d.Set("private_key", selfPrivateKey)
	} else {
		_ = d.Set("private_key", iam.FixPEM(createdService.PrivateKey))
	}

	// Set scopes and default_scopes
	_, _, err = client.Services.AddScopes(*createdService, scopes, defaultScopes)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if len(diags) > 0 {
		_, _, _ = client.Services.DeleteService(*createdService) // Cleanup
		return diags
	}

	d.SetId(createdService.ID)
	return resourceIAMServiceRead(ctx, d, m)
}

func resourceIAMServiceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	s, resp, err := client.Services.GetServiceByID(id)
	if err != nil {
		if err == iam.ErrEmptyResults || (resp != nil && resp.StatusCode == http.StatusNotFound) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("description", s.Description)
	_ = d.Set("name", s.Name)
	_ = d.Set("application_id", s.ApplicationID)
	_ = d.Set("organization_id", s.OrganizationID)
	_ = d.Set("service_id", s.ServiceID)
	_ = d.Set("scopes", s.Scopes)
	_ = d.Set("expires_on", s.ExpiresOn)
	_ = d.Set("default_scopes", s.DefaultScopes)
	return diags
}

func resourceIAMServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var s iam.Service
	s.ID = d.Id()
	s.ServiceID = d.Get("service_id").(string)

	if d.HasChange("scopes") {
		o, n := d.GetChange("scopes")
		old := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, old)
		toRemove := difference(old, newList)
		if len(toRemove) > 0 {
			_, _, err := client.Services.RemoveScopes(s, toRemove, []string{})
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if len(toAdd) > 0 {
			_, _, _ = client.Services.AddScopes(s, toAdd, []string{})
		}
	}
	if d.HasChange("default_scopes") {
		o, n := d.GetChange("default_scopes")
		old := expandStringList(o.(*schema.Set).List())
		newList := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newList, old)
		toRemove := difference(old, newList)
		if len(toRemove) > 0 {
			_, _, err := client.Services.RemoveScopes(s, []string{}, toRemove)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if len(toAdd) > 0 {
			_, _, _ = client.Services.AddScopes(s, []string{}, toAdd)
		}
	}
	if d.HasChange("private_key") || d.HasChange("expires_on") ||
		d.HasChange("self_managed_private_key") || d.HasChange("self_managed_certificate") {
		_, npk := d.GetChange("self_managed_private_key")
		_, npc := d.GetChange("self_managed_certificate")

		if npk.(string) == "" && npc.(string) == "" {
			return diag.FromErr(fmt.Errorf("you cannot revert to a server side managed private key once you set a self managed key or certificate"))
		}
		diags = setSelfManaged(client, s, d)
		if len(diags) > 0 {
			return diags
		}
	}

	return resourceIAMServiceRead(ctx, d, m)
}

func resourceIAMServiceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var s iam.Service
	s.ID = d.Id()
	ok, _, err := client.Services.DeleteService(s)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ok {
		return diag.FromErr(ErrDeleteServiceFailed)
	}
	d.SetId("")
	return diags
}

func setSelfManaged(client *iam.Client, service iam.Service, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	selfCertificate := d.Get("self_managed_certificate").(string)
	if selfCertificate != "" {
		block, _ := pem.Decode([]byte(iam.FixPEM(selfCertificate)))
		if block == nil {
			return diag.FromErr(fmt.Errorf("error decoding 'self_managed_certificate'"))
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing certificate: %w", err))
		}
		_, _, err = client.Services.UpdateServiceCertificateDER(service, cert.Raw)
		if err != nil {
			return diag.FromErr(fmt.Errorf("setting RAW certificate: %w", err))
		}
		return diags
	}
	selfPrivateKey := d.Get("self_managed_private_key").(string)
	selfExpiresOn := d.Get("expires_on").(string)
	expiresOn := time.Now().Add(5 * 86400 * 365 * time.Second)
	if selfExpiresOn != "" {
		parsedExpiresOn, err := time.Parse(time.RFC3339, selfExpiresOn)
		if err != nil {
			return diag.FromErr(fmt.Errorf("parsing expires_on: %w", err))
		}
		expiresOn = parsedExpiresOn
	}
	block, _ := pem.Decode([]byte(iam.FixPEM(selfPrivateKey)))
	if block == nil {
		block, _ = pem.Decode([]byte(selfPrivateKey)) // Try unmodified decode
		if block == nil {
			return diag.FromErr(fmt.Errorf("error decoding 'self_managed_private_key'"))
		}
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return diag.FromErr(fmt.Errorf("parsing private key: %w", err))
	}
	_, _, err = client.Services.UpdateServiceCertificate(service, privateKey, func(cert *x509.Certificate) error {
		cert.NotAfter = expiresOn
		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("setting private key: %w", err))
	}
	return diags
}
