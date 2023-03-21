package ch

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/loafoe/easyssh-proxy/v2"
	"github.com/philips-software/go-hsdp-api/cartel"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"

	"log"
	"net/http"
	"strings"
	"time"
)

const (
	fileField                  = "file"
	commandsField              = "commands"
	commandsDepecrationMessage = "The 'commands' argument is deprecated and will be removed in v0.40.0+. Please use the 'ssh_resource' from provider 'loafoe/ssh' for bootstrapping"
	fileDepecrationMessage     = "The 'file' block is deprecated and will be removed in v0.40.0+. Please use the 'ssh_resource' from provider 'loafoe/ssh' for bootstrapping"
)

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeMap,
		Required:         true,
		ValidateDiagFunc: validateTags,
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

func validateTags(v interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	tagsMap, ok := v.(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("expected %q to be a map", v))
	}
	if len(tagsMap) > 8 {
		return diag.FromErr(fmt.Errorf("maximum of 8 tags are supported"))
	}
	for k, v := range tagsMap {
		if strings.EqualFold(k, "name") {
			return diag.FromErr(fmt.Errorf("tag \"%s\" is reserved by the Cartel API", k))
		}
		val, ok := v.(string)
		if !ok {
			return diag.FromErr(fmt.Errorf("tag \"%s\" value is of type %q", k, v))
		}
		if len(val) > 255 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Detail:   fmt.Sprintf("value of tag \"%s\" is too long (max=255)", k),
			})
		}
	}
	return diags
}

func ResourceContainerHost() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceContainerHostCreate,
		ReadContext:   resourceContainerHostRead,
		UpdateContext: resourceContainerHostUpdate,
		DeleteContext: resourceContainerHostDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(25 * time.Minute),
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
			"image": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: validation.IntBetween(0, 16000),
			},
			"security_groups": {
				Type:     schema.TypeSet,
				MaxItems: 4,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user_groups": {
				Type:     schema.TypeSet,
				MaxItems: 50,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bastion_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"agent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"keep_failed_instances": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"commands_after_file_changes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			commandsField: {
				Type:       schema.TypeList,
				MaxItems:   10,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: commandsDepecrationMessage,
			},
			fileField: {
				Type:       schema.TypeSet,
				Optional:   true,
				Elem:       fileFieldSchema(),
				Deprecated: fileDepecrationMessage,
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
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
		SchemaVersion: 5,
	}
}

func fileFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func instanceStateRefreshFunc(client *cartel.Client, nameTag string, failStates []string) resource.StateRefreshFunc {
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
	c := m.(*config.Config)
	client, err := c.CartelClient()
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
	securityGroups := tools.ExpandStringList(d.Get("security_groups").(*schema.Set).List())
	userGroups := tools.ExpandStringList(d.Get("user_groups").(*schema.Set).List())
	instanceRole := d.Get("instance_role").(string)
	subnetType := d.Get("subnet_type").(string)
	bastionHost := d.Get("bastion_host").(string)
	keepFailedInstances := d.Get("keep_failed_instances").(bool)
	if bastionHost == "" {
		bastionHost = client.BastionHost()
	}
	image := d.Get("image").(string)
	user := d.Get("user").(string)
	privateKey := d.Get("private_key").(string)
	agent := d.Get("agent").(bool)

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
	// Validation
	if diags := validateContainerHostSchema(d); len(diags) > 0 {
		return diags
	}

	// Fetch files first before starting provisioning
	createFiles, diags := collectFilesToCreate(d)
	if len(diags) > 0 {
		return diags
	}
	// And commands
	commands, diags := tools.CollectList(commandsField, d)
	if len(diags) > 0 {
		return diags
	}

	if len(commands) > 0 || len(createFiles) > 0 {
		if user == "" && !agent {
			return diag.FromErr(fmt.Errorf("'user' must be set when 'agent = false' and '%s' are set or 'file' blocks are present", commandsField))
		}
		if privateKey == "" && !agent {
			return diag.FromErr(fmt.Errorf("no SSH 'private_key' was set and 'agent = false', authentication will fail after provisioning step"))
		}
		if agent && !tools.SSHAgentReachable() {
			return diag.FromErr(fmt.Errorf("'agent = true' but no working 'ssh-agent' socket is advertised in SSH_AUTH_SOCK environment variable"))
		}
	}

	instanceID := ""
	ipAddress := ""
	needCreate := true
	// First, check if the instance already exists for whatever reason. Cartel / AWS are flaky sometimes
	if details := findInstanceByName(client, tagName); details != nil {
		testTag := uuid.NewString()
		// Next, try to set and unset a tag, to ensure we can own/control it
		_, _, err := client.AddTags([]string{tagName}, map[string]string{
			"tf-crud-check": testTag,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("no write access to instance '%s', giving up", tagName))
		}
		if details.Tags == nil {
			details.Tags = make(map[string]string)
		}
		details.Tags["tf-crud-check"] = ""
		_, _, _ = client.AddTags([]string{tagName}, details.Tags)
		needCreate = false
		instanceID = details.InstanceID
		ipAddress = details.PrivateAddress
	}

	if needCreate {
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
			cartel.Image(image),
		)
		if err != nil {
			// Do not clean up existing hosts
			if err == cartel.ErrHostnameAlreadyExists {
				return diag.FromErr(fmt.Errorf("the host '%s' already exists: %w", tagName, err))
			}
			if resp == nil {
				if keepFailedInstances {
					diags = append(diags, diag.FromErr(fmt.Errorf("'keep_failed_instances' is enabled so not removing '%s', remember to destroy it manually", tagName))...)
				} else {
					_, _, _ = client.Destroy(tagName)
				}
				diags = append(diags, diag.FromErr(fmt.Errorf("create error (resp=nil): %w", err))...)
				return diags
			}
			if ch == nil || resp.StatusCode() >= 500 { // Possible 504, or other timeout, try to recover!
				if details := findInstanceByName(client, tagName); details != nil {
					instanceID = details.InstanceID
					ipAddress = details.PrivateAddress
				} else {
					if keepFailedInstances {
						diags = append(diags, diag.FromErr(fmt.Errorf("'keep_failed_instances' is enabled so not removing '%s', remember to destroy it manually", tagName))...)
					} else {
						_, _, _ = client.Destroy(tagName)
					}
					diags = append(diags, diag.FromErr(fmt.Errorf("create error (status=%d): %w", resp.StatusCode(), err))...)
					return diags
				}
			} else {
				if keepFailedInstances {
					diags = append(diags, diag.FromErr(fmt.Errorf("'keep_failed_instances' is enabled so not removing '%s', remember to destroy it manually", tagName))...)
				} else {
					_, _, _ = client.Destroy(tagName)
				}
				diags = append(diags, diag.FromErr(fmt.Errorf("create error (description=[%s], code=[%d]): %w", ch.Description, resp.StatusCode(), err))...)
				return diags
			}
		} else {
			instanceID = ch.InstanceID()
			ipAddress = ch.IPAddress()
		}
	}

	// Randomize checks a bit in case of many concurrent creates
	rand.Seed(time.Now().UnixNano())
	min := 5
	max := 15
	minTimeout := rand.Intn(max-min+1) + min

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"provisioning", "indeterminate"},
		Target:     []string{"succeeded"},
		Refresh:    instanceStateRefreshFunc(client, tagName, []string{"failed", "terminated", "shutting-down"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: time.Duration(minTimeout) * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		// Trigger a delete to prevent failed instances from lingering
		if !keepFailedInstances {
			_, _, _ = client.Destroy(tagName)
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf(
			"error waiting for instance '%s' to become ready: %v",
			instanceID, err))
	}
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": ipAddress,
	})
	// Collect SSH details
	privateIP := ipAddress
	ssh := &easyssh.MakeConfig{
		User:   user,
		Server: privateIP,
		Port:   "22",
		Proxy:  http.ProxyFromEnvironment,
		Bastion: easyssh.DefaultConfig{
			User:   user,
			Server: bastionHost,
			Port:   "22",
		},
	}
	if privateKey != "" {
		ssh.Key = privateKey
		ssh.Bastion.Key = privateKey
	}

	// Check health of Docker daemon in case of 'container-host' role and file or commands are set
	if (len(commands) > 0 || len(createFiles) > 0) && instanceRole == "container-host" {
		if err := ensureContainerHostReady(ssh, c); err != nil {
			if !keepFailedInstances {
				_, _, _ = client.Destroy(tagName)
				d.SetId("")
			}
			return diag.FromErr(fmt.Errorf(
				"container host instance '%s' was not deemed healthy: %v",
				instanceID, err))
		}
	}

	// Create files
	_, _ = c.Debug("about to copy %d files to remote\n", len(createFiles))
	if err := copyFiles(ssh, c, createFiles); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "failed to copy all files",
			Detail:   fmt.Sprintf("One or more files failed to copy: %v", err),
		})
	}

	// Run commands
	stdout, errDiags, err := runCommands(commands, ssh, m)
	if err != nil {
		return errDiags
	}
	_ = d.Set("result", stdout)
	d.SetId(instanceID)
	readDiags := resourceContainerHostRead(ctx, d, m)
	return append(diags, readDiags...)
}

func ensureContainerHostReady(ssh *easyssh.MakeConfig, config *config.Config) error {
	operation := func() error {
		outStr, errStr, done, err := ssh.Run("docker volume ls") // This command should succeed
		_, _ = config.Debug("ensureContainerHostReady: %t\nstdout:\n%s\nstderr:\n%s\n", done, outStr, errStr)
		if err != nil {
			if err.Error() == "Process exited with status 1" { // Currently, the only known value to retry on
				return err
			}
			return backoff.Permanent(err)
		}
		return nil
	}
	err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 20))
	if err != nil {
		return err
	}
	return nil
}

func validateContainerHostSchema(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	securityGroups := tools.ExpandStringList(d.Get("security_groups").(*schema.Set).List())

	if tools.ContainsString(securityGroups, "base") {
		return diag.FromErr(fmt.Errorf("the 'base' security group is internal and should not be specified"))
	}
	return diags
}

func findInstanceByName(client *cartel.Client, name string) *cartel.InstanceDetails {
	instances, _, err := client.GetAllInstances()
	if err != nil {
		return nil
	}
	for _, i := range *instances {
		if i.NameTag == name {
			return &i
		}
	}
	return nil
}

func copyFiles(ssh *easyssh.MakeConfig, config *config.Config, createFiles []provisionFile) error {
	for _, f := range createFiles {
		if f.Source != "" {
			src, srcErr := os.Open(f.Source)
			if srcErr != nil {
				_, _ = config.Debug("Failed to open source file %s: %v\n", f.Source, srcErr)
				return srcErr
			}
			srcStat, statErr := src.Stat()
			if statErr != nil {
				_, _ = config.Debug("Failed to stat source file %s: %v\n", f.Source, statErr)
				_ = src.Close()
				return fmt.Errorf("copyFiles: %w", statErr)
			}
			err := ssh.WriteFile(src, srcStat.Size(), f.Destination)
			if err != nil {
				_, _ = config.Debug("Error copying %s to remote file %s:%s: %v\n", f.Source, ssh.Server, f.Destination, err)
				return fmt.Errorf("copyFiles: %w", err)
			}
			_, _ = config.Debug("Copied %s to remote file %s:%s: %d bytes\n", f.Source, ssh.Server, f.Destination, srcStat.Size())
			_ = src.Close()
		} else {
			buffer := bytes.NewBufferString(f.Content)
			// Should we fail the complete provision on errors here?
			err := ssh.WriteFile(buffer, int64(buffer.Len()), f.Destination)
			if err != nil {
				_, _ = config.Debug("Error copying content to remote file %s:%s: %v\n", ssh.Server, f.Destination, err)
				return fmt.Errorf("copyFiles: %w", err)
			}
			_, _ = config.Debug("Created remote file %s:%s: %d bytes\n", ssh.Server, f.Destination, len(f.Content))
		}
		// Permissions change
		if f.Permissions != "" {
			outStr, errStr, _, err := ssh.Run(fmt.Sprintf("chmod %s \"%s\"", f.Permissions, f.Destination))
			_, _ = config.Debug("Permissions file %s:%s: %v %v\n", f.Destination, f.Permissions, outStr, errStr)
			if err != nil {
				return err
				// Owner
			}
		}
		if f.Owner != "" {
			outStr, errStr, _, err := ssh.Run(fmt.Sprintf("chown %s \"%s\"", f.Owner, f.Destination))
			_, _ = config.Debug("Owner file %s:%s: %v %v\n", f.Destination, f.Owner, outStr, errStr)
			if err != nil {
				return err
			}
		}
		// Group
		if f.Group != "" {
			outStr, errStr, _, err := ssh.Run(fmt.Sprintf("chgrp %s \"%s\"", f.Group, f.Destination))
			_, _ = config.Debug("Group file %s:%s: %v %v\n", f.Destination, f.Group, outStr, errStr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type provisionFile struct {
	Source      string
	Content     string
	Destination string
	Permissions string
	Owner       string
	Group       string
}

func collectFilesToCreate(d *schema.ResourceData) ([]provisionFile, diag.Diagnostics) {
	var diags diag.Diagnostics
	files := make([]provisionFile, 0)
	if v, ok := d.GetOk(fileField); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			file := provisionFile{
				Source:      mVi["source"].(string),
				Content:     mVi["content"].(string),
				Destination: mVi["destination"].(string),
				Permissions: mVi["permissions"].(string),
				Owner:       mVi["owner"].(string),
				Group:       mVi["group"].(string),
			}
			if file.Source == "" && file.Content == "" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "conflict in file block",
					Detail:   fmt.Sprintf("file %s has neither 'source' or 'content', set one", file.Destination),
				})
				continue
			}
			if file.Source != "" && file.Content != "" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "conflict in file block",
					Detail:   fmt.Sprintf("file %s has conflicting 'source' and 'content', choose only one", file.Destination),
				})
				continue
			}
			if file.Source != "" {
				src, srcErr := os.Open(file.Source)
				if srcErr != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "issue with source",
						Detail:   fmt.Sprintf("file %s: %v", file.Source, srcErr),
					})
					continue
				}
				_, statErr := src.Stat()
				if statErr != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "issue with source stat",
						Detail:   fmt.Sprintf("file %s: %v", file.Source, statErr),
					})
					_ = src.Close()
					continue
				}
				_ = src.Close()
			}
			files = append(files, file)
		}
	}
	return files, diags
}

func resourceContainerHostUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagName := d.Get("name").(string)
	ch, _, err := client.GetDetails(tagName)
	if err != nil {
		return diag.FromErr(err)
	}
	if ch.InstanceID != d.Id() {
		return diag.FromErr(config.ErrInstanceIDMismatch)
	}

	// Validation
	if diags := validateContainerHostSchema(d); len(diags) > 0 {
		return diags
	}
	bastionHost := d.Get("bastion_host").(string)
	user := d.Get("user").(string)
	privateKey := d.Get("private_key").(string)
	commandsAfterFileChanges := d.Get("commands_after_file_changes").(bool)
	agent := d.Get("agent").(bool)
	if bastionHost == "" {
		bastionHost = client.BastionHost()
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
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newEntries := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newEntries, old)
		toRemove := tools.Difference(old, newEntries)

		// Removals
		if len(toRemove) > 0 {
			_, _, err := client.RemoveUserGroups([]string{tagName}, toRemove)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Additions
		if len(toAdd) > 0 {
			_, _, err := client.AddUserGroups([]string{tagName}, toAdd)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("security_groups") {
		o, n := d.GetChange("security_groups")
		old := tools.ExpandStringList(o.(*schema.Set).List())
		newEntries := tools.ExpandStringList(n.(*schema.Set).List())
		toAdd := tools.Difference(newEntries, old)
		toRemove := tools.Difference(old, newEntries)

		// Removals
		if len(toRemove) > 0 {
			_, _, err := client.RemoveSecurityGroups([]string{tagName}, toRemove)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Additions
		if len(toAdd) > 0 {
			_, _, err := client.AddSecurityGroups([]string{tagName}, toAdd)
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
	// Collect SSH details
	privateIP := d.Get("private_ip").(string)
	ssh := &easyssh.MakeConfig{
		User:   user,
		Server: privateIP,
		Port:   "22",
		Proxy:  http.ProxyFromEnvironment,
		Bastion: easyssh.DefaultConfig{
			User:   user,
			Server: bastionHost,
			Port:   "22",
		},
	}
	if privateKey != "" {
		if agent {
			return diag.FromErr(fmt.Errorf("'agent' is enabled so not expecting a private key to be set"))
		}
		ssh.Key = privateKey
		ssh.Bastion.Key = privateKey
	}
	if d.HasChange("file") {
		createFiles, diags := collectFilesToCreate(d)
		if len(diags) > 0 {
			return diags
		}
		_, _ = c.Debug("about to copy %d files to remote\n", len(createFiles))
		if err := copyFiles(ssh, c, createFiles); err != nil {
			return diag.FromErr(fmt.Errorf("copying files to remote: %w", err))
		}
		if commandsAfterFileChanges {
			commands, diags := tools.CollectList(commandsField, d)
			if len(diags) > 0 {
				return diags
			}
			// Run commands
			stdout, errDiags, err := runCommands(commands, ssh, m)
			if err != nil {
				return errDiags
			}
			_ = d.Set("result", stdout)
		}
	}
	return diags
}

func resourceContainerHostRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagName := d.Get("name").(string)

	if tagName == "" { // This is an import, find and set the tagName
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
		if resp != nil && resp.StatusCode() == http.StatusBadRequest {
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
		return diag.FromErr(config.ErrInstanceIDMismatch)
	}
	_ = d.Set("protect", ch.Protection)
	_ = d.Set("volumes", len(ch.BlockDevices)-1) // -1 for the root volume
	_ = d.Set("role", ch.Role)
	_ = d.Set("launch_time", ch.LaunchTime)
	_ = d.Set("block_devices", ch.BlockDevices)
	_ = d.Set("security_groups", tools.Difference(ch.SecurityGroups, []string{"base"})) // Remove "base"
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
	c := m.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagName := d.Get("name").(string)

	ch, _, err := client.GetDetails(tagName)
	if err != nil {
		return diag.FromErr(err)
	}
	if ch.InstanceID != d.Id() {
		return diag.FromErr(config.ErrInstanceIDMismatch)
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

func runCommands(commands []string, ssh *easyssh.MakeConfig, m interface{}) (string, diag.Diagnostics, error) {
	var diags diag.Diagnostics
	var stdout, stderr string
	var done bool
	var err error
	c := m.(*config.Config)

	for i := 0; i < len(commands); i++ {
		stdout, stderr, done, err = ssh.Run(commands[i], 5*time.Minute)
		_, _ = c.Debug("command: %s\ndone: %t\nstdout:\n%s\nstderr:\n%s\n", commands[i], done, stdout, stderr)
		if err != nil {
			_, _ = c.Debug("error: %v\n", err)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("execution of command '%s' failed. stdout output", commands[i]),
				Detail:   stdout,
			})
			if stderr != "" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "stderr output",
					Detail:   stderr,
				})
			}
			return stdout, diags, err
		}
	}
	return stdout, diags, nil
}
