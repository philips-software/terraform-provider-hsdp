package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/pkg/errors"
)

func ResourceIAMActivationEmail() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMActivationEmailCreate,
		ReadContext:   resourceIAMActivationEmailRead,
		UpdateContext: resourceIAMActivationEmailUpdate,
		DeleteContext: resourceIAMActivationEmailDelete,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resend_every": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  72,
				ForceNew: true,
			},
			"send": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"last_sent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Second),
			Update: schema.DefaultTimeout(300 * time.Second),
		},
	}
}

func resourceIAMActivationEmailUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("send") {
		loginID := d.Get("login_id").(string)
		_, _, err := client.Users.ResendActivation(loginID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error sending activation: %w", err))
		}
		_ = d.Set("last_sent", time.Now().UTC().Format(time.RFC3339))
	}
	return diags
}

func resourceIAMActivationEmailRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	userID := d.Get("user_id").(string)
	sendEvery := d.Get("resend_every").(int)

	user, _, err := client.Users.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, iam.ErrEmptyResults) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("read user: %w", err))
	}

	_ = d.Set("verified", user.AccountStatus.EmailVerified)
	_ = d.Set("login_id", user.LoginID)
	lastSent, err := time.Parse(time.RFC3339, d.Get("last_sent").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("cannot determine last_sent: %w", err))
	}
	nextSend := lastSent.Add(time.Duration(sendEvery) * time.Hour)
	if time.Until(nextSend) <= 0 { // Time to send
		_ = d.Set("send", true)
	} else {
		_ = d.Set("send", nil)
	}

	return diags
}

func resourceIAMActivationEmailCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get("user_id").(string)

	user, _, err := client.Users.GetUserByID(userID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("read user: %w", err))
	}
	_ = d.Set("verified", user.AccountStatus.EmailVerified)
	_ = d.Set("last_sent", time.Now().UTC().Format(time.RFC3339))
	_ = d.Set("email_address", user.EmailAddress)

	d.SetId(userID)
	return diags
}

func resourceIAMActivationEmailDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
