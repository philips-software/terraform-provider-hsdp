package hsdp

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	creds "github.com/philips-software/go-hsdp-api/credentials"
)

func resourceCredentialsPolicy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceCredentialsPolicyCreate,
		Read:   resourceCredentialsPolicyRead,
		Delete: resourceCredentialsPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validatePolicyJSON,
				DiffSuppressFunc: suppressEquivalentPolicyDiffs,
			},
			"product_key": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				ForceNew:  true,
				Required:  true,
			},
		},
	}
}

func resourceCredentialsPolicyCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.CredentialsClient()
	if err != nil {
		return err
	}

	productKey := d.Get("product_key").(string)
	policyJSON := d.Get("policy").(string)
	var policy creds.Policy

	err = json.Unmarshal([]byte(policyJSON), &policy)
	if err != nil {
		return err
	}
	policy.ProductKey = productKey

	createdPolicy, _, err := client.Policy.CreatePolicy(policy)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(createdPolicy.ID))
	return nil
}

func resourceCredentialsPolicyRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.CredentialsClient()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	productKey := d.Get("product_key").(string)

	policies, _, err := client.Policy.GetPolicy(&creds.GetPolicyOptions{
		ID:         &id,
		ProductKey: &productKey,
	})
	if err != nil {
		return err
	}
	if len(policies) != 1 { // Policy was deleted
		d.SetId("")
		return nil
	}
	policy := policies[0]

	d.SetId(strconv.Itoa(policy.ID))
	policy.ID = 0 // Don't marshal ID
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("policy", policyJSON)
	d.Set("product_key", productKey)
	return nil
}

func resourceCredentialsPolicyDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client, err := config.CredentialsClient()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	productKey := d.Get("product_key").(string)
	policy := creds.Policy{
		ID:         id,
		ProductKey: productKey,
	}
	ok, _, err := client.Policy.DeletePolicy(policy)
	if err != nil {
		return err
	}
	if ok {
		d.SetId("")
	}
	return nil
}
