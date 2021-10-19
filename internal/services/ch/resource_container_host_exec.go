package ch

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/loafoe/easyssh-proxy/v2"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
)

func ResourceContainerHostExec() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`hsdp_container_host_exec`" + ` resource implements the standard resource lifecycle but takes no further action.
The ` + "`triggers`" + ` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.`,

		CreateContext: resourceContainerHostExecCreate,
		Read:          resourceContainerHostExecRead,
		Delete:        resourceContainerHostExecDelete,
		SchemaVersion: 2,

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
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"agent": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			commandsField: {
				Type:     schema.TypeList,
				MaxItems: 50,
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
	c := m.(*config.Config)
	client, err := c.CartelClient()
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
	agent := d.Get("agent").(bool)

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
		if user == "" {
			return diag.FromErr(fmt.Errorf("user must be set when '%s' is specified", commandsField))
		}
		if privateKey == "" && !agent {
			return diag.FromErr(fmt.Errorf("no SSH 'private_key' was set and 'agent' is 'false', authentication will fail after provisioning step"))
		}
	}
	// Collect SSH details
	privateIP := host
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

	// Provision files
	if err := copyFiles(ssh, c, createFiles); err != nil {
		return diag.FromErr(fmt.Errorf("copying files to remote: %w", err))
	}

	// Ensure ready-ness
	if err := ensureContainerHostReady(ssh, c); err != nil {
		return diag.FromErr(fmt.Errorf("container host ready-ness check failed: %w", err))
	}

	// Run commands
	stdout, errDiags, err := runCommands(commands, ssh, m)
	if err != nil {
		return errDiags
	}
	_ = d.Set("result", stdout)
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return diags
}

func resourceContainerHostExecRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceContainerHostExecDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
