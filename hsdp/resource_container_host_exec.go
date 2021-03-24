package hsdp

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/loafoe/easyssh-proxy/v2"
)

func resourceContainerHostExec() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`hsdp_container_host_exec`" + ` resource implements the standard resource lifecycle but takes no further action.
The ` + "`triggers`" + ` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.`,

		CreateContext: resourceContainerHostExecCreate,
		Read:          resourceContainerHostExecRead,
		Delete:        resourceConainterHostDelete,

		Schema: map[string]*schema.Schema{
			"triggers": {
				Description: "A map of arbitrary strings that, when changed, will force the 'hsdp_container_host_exec' resource to be replaced, re-running any associated commands.",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bastion_host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			commandsField: {
				Type:     schema.TypeList,
				MaxItems: 10,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			fileField: {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
				},
			},
		},
	}
}

func resourceContainerHostExecCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client, err := config.CartelClient()
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	bastionHost := d.Get("bastion_host").(string)
	if bastionHost == "" {
		bastionHost = client.BastionHost()
	}
	user := d.Get("user").(string)
	privateKey := d.Get("private_key").(string)
	host := d.Get("host").(string)

	// Fetch files first before starting provisioning
	createFiles, diags := collectFilesToCreate(d)
	if len(diags) > 0 {
		return diags
	}
	// And commands
	commands, diags := collectCommands(d)
	if len(diags) > 0 {
		return diags
	}
	if len(commands) > 0 {
		if user == "" {
			return diag.FromErr(fmt.Errorf("user must be set when '%s' is specified", commandsField))
		}
		if privateKey == "" {
			return diag.FromErr(fmt.Errorf("privateKey must be set when '%s' is specified", commandsField))
		}
	}
	// Collect SSH details
	privateIP := host
	ssh := &easyssh.MakeConfig{
		User:   user,
		Server: privateIP,
		Port:   "22",
		Key:    privateKey,
		Proxy:  http.ProxyFromEnvironment,
		Bastion: easyssh.DefaultConfig{
			User:   user,
			Server: bastionHost,
			Port:   "22",
			Key:    privateKey,
		},
	}

	// Provision files
	if err := copyFiles(ssh, config, createFiles); err != nil {
		return diag.FromErr(fmt.Errorf("copying files to remote: %w", err))
	}

	// Run commands
	for i := 0; i < len(commands); i++ {
		stdout, stderr, done, err := ssh.Run(commands[i], 5*time.Minute)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("command [%s]: %w", commands[i], err))...)
		} else {
			_, _ = config.Debug("command: %s\ndone: %t\nstdout:\n%s\nstderr:\n%s\n", commands[i], done, stdout, stderr)
		}
	}

	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return diags
}

func resourceContainerHostExecRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceConainterHostDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
