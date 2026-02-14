package iam

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-dip-api/iam"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceIAMPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceIAMPasswordPolicyCreate,
		UpdateContext: resourceIAMPasswordPolicyUpdate,
		ReadContext:   resourceIAMPasswordPolicyRead,
		DeleteContext: resourceIAMPasswordPolicyDelete,

		Schema: map[string]*schema.Schema{
			"managing_organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expiry_period_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  90,
			},
			"history_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"complexity": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  8,
						},
						"max_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  16,
						},
						"min_numerics": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"min_uppercase": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"min_lowercase": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"min_special_chars": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
			"challenges_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"challenge_policy": {
				Type:         schema.TypeSet,
				Optional:     true,
				RequiredWith: []string{"challenges_enabled"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_questions": {
							Type:     schema.TypeSet,
							MaxItems: 10,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"min_question_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"min_answer_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"max_incorrect_attempts": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
					},
				},
			},
			"_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDataToToPasswordPolicy(d *schema.ResourceData, policy *iam.PasswordPolicy) {
	policy.ManagingOrganization = d.Get("managing_organization").(string)
	policy.ExpiryPeriodInDays = d.Get("expiry_period_in_days").(int)
	policy.HistoryCount = d.Get("history_count").(int)
	if v, ok := d.GetOk("complexity"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			policy.Complexity.MaxLength = mVi["max_length"].(int)
			policy.Complexity.MinLength = mVi["min_length"].(int)
			policy.Complexity.MinLowerCase = mVi["min_lowercase"].(int)
			policy.Complexity.MinUpperCase = mVi["min_uppercase"].(int)
			policy.Complexity.MinNumerics = mVi["min_numerics"].(int)
			policy.Complexity.MinSpecialChars = mVi["min_special_chars"].(int)
		}
	}
	policy.ChallengesEnabled = d.Get("challenges_enabled").(bool)
	if v, ok := d.GetOk("challenge_policy"); ok {
		policy.ChallengePolicy = &iam.ChallengePolicy{}
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			policy.ChallengePolicy.DefaultQuestions = tools.ExpandStringList(mVi["default_questions"].(*schema.Set).List())
			policy.ChallengePolicy.MinQuestionCount = mVi["min_question_count"].(int)
			policy.ChallengePolicy.MinAnswerCount = mVi["min_answer_count"].(int)
			policy.ChallengePolicy.MaxIncorrectAttempts = mVi["max_incorrect_attempts"].(int)
		}
	}
}

func resourceIAMPasswordPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var policy iam.PasswordPolicy
	resourceDataToToPasswordPolicy(d, &policy)

	// Since there's only a single password policy per ORG, first try to fetch it
	policies, _, err := client.PasswordPolicies.GetPasswordPolicies(&iam.GetPasswordPolicyOptions{
		OrganizationID: &policy.ManagingOrganization,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	policyFunc := client.PasswordPolicies.CreatePasswordPolicy
	if policies != nil && len(*policies) > 0 {
		existingPolicy := (*policies)[0]
		policy.ID = existingPolicy.ID
		policy.Meta = existingPolicy.Meta
		policyFunc = client.PasswordPolicies.UpdatePasswordPolicy
	}
	var createdPolicy *iam.PasswordPolicy

	err = tools.TryHTTPCall(ctx, 10, func() (*http.Response, error) {
		var err error
		var resp *iam.Response
		createdPolicy, resp, err = policyFunc(policy)
		if err != nil {
			_ = client.TokenRefresh()
		}
		if resp == nil {
			return nil, err
		}
		return resp.Response, err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	data, err := json.Marshal(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("_policy", string(data))
	d.SetId(createdPolicy.ID)
	return diags
}

func resourceIAMPasswordPolicyUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}
	var updatePolicy iam.PasswordPolicy

	policy, _, err := client.PasswordPolicies.GetPasswordPolicyByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resourceDataToToPasswordPolicy(d, &updatePolicy)
	updatePolicy.ID = policy.ID
	updatePolicy.Meta = policy.Meta

	updatedPolicy, _, err := client.PasswordPolicies.UpdatePasswordPolicy(updatePolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	data, err := json.Marshal(updatedPolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("_policy", string(data))
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourceIAMPasswordPolicyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policy, _, err := client.PasswordPolicies.GetPasswordPolicyByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	data, err := json.Marshal(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("_policy", data)
	d.SetId(policy.ID)
	return diags
}

func resourceIAMPasswordPolicyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var policy iam.PasswordPolicy
	policy.ID = d.Id()
	ok, _, err := client.PasswordPolicies.DeletePasswordPolicy(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	if ok {
		d.SetId("")
	}
	return diags
}
