package hsdp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-software/go-hsdp-api/cartel"
	"log"
	"net/http"
	"time"
)

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Required: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// TODO: handle empty tags
			return k == "tags.billing"
		},
		DefaultFunc: func() (interface{}, error) {
			return map[string]interface{}{"billing": ""}, nil
		},
		Elem: &schema.Schema{Type: schema.TypeString},
	}
}

func resourceContainerHost() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceContainerHostCreate,
		ReadContext:   resourceContainerHostRead,
		UpdateContext: resourceContainerHostUpdate,
		DeleteContext: resourceContainerHostDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "container-host",
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "m5.large",
			},
			"volume_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"iops"},
			},
			"iops": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 4000),
			},
			"protect": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"encrypt_volumes": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
				ForceNew: true,
			},
			"volumes": {
				Type:         schema.TypeInt,
				Default:      0,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 6),
			},
			"volume_size": {
				Type:         schema.TypeInt,
				Default:      0,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"security_groups": {
				Type:     schema.TypeSet,
				MaxItems: 5,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user_groups": {
				Type:     schema.TypeSet,
				MaxItems: 50,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"subnet"},
			},
			"subnet": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"launch_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_devices": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
		},
		SchemaVersion: 2,
	}
}

func InstanceStateRefreshFunc(client *cartel.Client, nameTag string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		state, resp, err := client.GetDeploymentState(nameTag)
		if err != nil {
			log.Printf("Error on InstanceStateRefresh: %s", err)
			return resp, "", err
		}

		for _, failState := range failStates {
			if state == failState {
				return resp, state, fmt.Errorf("failed to reach target state, reason: %s",
					state)
			}
		}
		return resp, state, nil
	}
}

func resourceContainerHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client, err := config.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagName := d.Get("name").(string)
	protect := d.Get("protect").(bool)
	iops := d.Get("iops").(int)
	encryptVolumes := d.Get("encrypt_volumes").(bool)
	volumeSize := d.Get("volume_size").(int)
	numberOfVolumes := d.Get("volumes").(int)
	volumeType := d.Get("volume_type").(string)
	instanceType := d.Get("instance_type").(string)
	securityGroups := expandStringList(d.Get("security_groups").(*schema.Set).List())
	userGroups := expandStringList(d.Get("user_groups").(*schema.Set).List())
	instanceRole := d.Get("instance_role").(string)
	subnetType := d.Get("subnet_type").(string)
	if subnetType == "" {
		subnetType = "private"
	}
	subnet := d.Get("subnet").(string)
	tagList := d.Get("tags").(map[string]interface{})
	tags := make(map[string]string)
	for t, v := range tagList {
		if val, ok := v.(string); ok {
			tags[t] = val
		}
	}

	ch, resp, err := client.Create(tagName,
		cartel.SecurityGroups(securityGroups...),
		cartel.UserGroups(userGroups...),
		cartel.VolumeType(volumeType),
		cartel.IOPs(iops),
		cartel.InstanceType(instanceType),
		cartel.VolumesAndSize(numberOfVolumes, volumeSize),
		cartel.VolumeEncryption(encryptVolumes),
		cartel.Protect(protect),
		cartel.InstanceRole(instanceRole),
		cartel.SubnetType(subnetType),
		cartel.Tags(tags),
		cartel.InSubnet(subnet),
	)
	if err != nil {
		if resp == nil {
			return diag.FromErr(fmt.Errorf("create error (resp=nil): %w", err))
		}
		if ch == nil {
			return diag.FromErr(fmt.Errorf("create error (instance=nil): %w", err))
		}
		return diag.FromErr(fmt.Errorf("create error (description=[%s]): %w", ch.Description, err))
	}
	d.SetId(ch.InstanceID())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"provisioning", "indeterminate"},
		Target:     []string{"succeeded"},
		Refresh:    InstanceStateRefreshFunc(client, tagName, []string{"failed", "terminated", "shutting-down"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		// Trigger a delete to prevent failed instances from lingering
		_, _, _ = client.Destroy(tagName)
		return diag.FromErr(fmt.Errorf(
			"error waiting for instance (%s) to become ready: %s",
			ch.InstanceID(), err))
	}
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": ch.IPAddress(),
	})
	return resourceContainerHostRead(ctx, d, m)
}

func resourceContainerHostUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagName := d.Get("name").(string)
	ch, _, err := client.GetDetails(tagName)
	if err != nil {
		return diag.FromErr(err)
	}
	if ch.InstanceID != d.Id() {
		return diag.FromErr(ErrInstanceIDMismatch)
	}

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		change := generateTagChange(o, n)
		log.Printf("[o:%v] [n:%v] [c:%v]\n", o, n, change)
		_, _, err := client.AddTags([]string{tagName}, change)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("user_groups") {
		o, n := d.GetChange("user_groups")
		old := expandStringList(o.(*schema.Set).List())
		newEntries := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newEntries, old)
		toRemove := difference(old, newEntries)

		// Additions
		if len(toAdd) > 0 {
			_, _, err := client.AddUserGroups([]string{tagName}, toAdd)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Removals
		if len(toRemove) > 0 {
			_, _, err := client.RemoveUserGroups([]string{tagName}, toRemove)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("security_groups") {
		o, n := d.GetChange("security_groups")
		old := expandStringList(o.(*schema.Set).List())
		newEntries := expandStringList(n.(*schema.Set).List())
		toAdd := difference(newEntries, old)
		toRemove := difference(old, newEntries)

		// Additions
		if len(toAdd) > 0 {
			_, _, err := client.AddSecurityGroups([]string{tagName}, toAdd)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Removals
		if len(toRemove) > 0 {
			_, _, err := client.RemoveSecurityGroups([]string{tagName}, toRemove)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("protect") {
		protect := d.Get("protect").(bool)
		_, _, err := client.SetProtection(tagName, protect)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags

}

func resourceContainerHostRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagName := d.Get("name").(string)

	if tagName == "" { // This an import, find and set the tagName
		instances, _, err := client.GetAllInstances()
		if err != nil {
			return diag.FromErr(fmt.Errorf("cartel.GetAllInstances: %w", err))
		}
		id := d.Id()
		_ = d.Set("encrypt_volumes", true)
		_ = d.Set("volume_size", 0)
		for _, i := range *instances {
			if i.InstanceID == id {
				_ = d.Set("name", i.NameTag)
				tagName = i.NameTag
				break
			}
		}
	}

	state, resp, err := client.GetDeploymentState(tagName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusBadRequest {
			// State not found, probably a botched provision :(
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	if state != "succeeded" {
		// Unless we have a succeeded deploy, taint the resource
		d.SetId("")
		return diags
	}
	ch, _, err := client.GetDetails(tagName)
	if err != nil {
		return diag.FromErr(err)
	}
	if ch.InstanceID != d.Id() {
		return diag.FromErr(ErrInstanceIDMismatch)
	}
	_ = d.Set("protect", ch.Protection)
	_ = d.Set("volumes", len(ch.BlockDevices)-1) // -1 for the root volume
	_ = d.Set("role", ch.Role)
	_ = d.Set("launch_time", ch.LaunchTime)
	_ = d.Set("block_devices", ch.BlockDevices)
	_ = d.Set("security_groups", difference(ch.SecurityGroups, []string{"base"})) // Remove "base"
	_ = d.Set("user_groups", ch.LdapGroups)
	_ = d.Set("instance_type", ch.InstanceType)
	_ = d.Set("instance_role", ch.Role)
	_ = d.Set("vpc", ch.Vpc)
	_ = d.Set("zone", ch.Zone)
	_ = d.Set("launch_time", ch.LaunchTime)
	_ = d.Set("private_ip", ch.PrivateAddress)
	_ = d.Set("public_ip", ch.PublicAddress)
	_ = d.Set("subnet", ch.Subnet)
	subnetType := "private"
	if ch.PublicAddress != "" {
		subnetType = "public"
	}
	_ = d.Set("subnet_type", subnetType)
	_ = d.Set("tags", normalizeTags(ch.Tags))

	return diags
}

func resourceContainerHostDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	client, err := config.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagName := d.Get("name").(string)
	ch, _, err := client.GetDetails(tagName)
	if err != nil {
		return diag.FromErr(err)
	}
	if ch.InstanceID != d.Id() {
		return diag.FromErr(ErrInstanceIDMismatch)
	}
	_, _, err = client.Destroy(tagName)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags

}

func normalizeTags(tags map[string]string) map[string]string {
	normalized := make(map[string]string)
	for k, v := range tags {
		if k == "billing" || v == "" {
			continue
		}
		normalized[k] = v
	}
	return normalized
}

func generateTagChange(old, new interface{}) map[string]string {
	change := make(map[string]string)
	o := old.(map[string]interface{})
	n := new.(map[string]interface{})
	for k := range o {
		if newVal, ok := n[k]; !ok || newVal == "" {
			change[k] = ""
		}
	}
	for k, v := range n {
		if k == "billing" {
			continue
		}
		if s, ok := v.(string); ok {
			change[k] = s
		}
	}
	return change
}
