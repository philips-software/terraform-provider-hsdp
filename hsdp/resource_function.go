package hsdp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/iron"
)

func resourceFunction() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CreateContext: resourceFunctionCreate,
		ReadContext:   resourceFunctionRead,
		UpdateContext: resourceFunctionUpdate,
		DeleteContext: resourceFunctionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"docker_image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"docker_credentials": {
				Type:      schema.TypeMap,
				Optional:  true,
				Sensitive: true,
			},
			"command": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"environment": {
				Type:      schema.TypeMap,
				ForceNew:  true,
				Optional:  true,
				Sensitive: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"run_every": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"backend": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"credentials": {
							Type:      schema.TypeMap,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
							//ValidateFunc: validation.StringIsJSON,
						},
					},
				},
			},
		},
	}
}

func resourceFunctionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ironClient, _, _, err := newIronClient(d)
	if err != nil {
		return diag.FromErr(err)
	}
	ids := strings.Split(d.Id(), "-")
	taskType := ids[0]
	taskID := ids[1]
	codeID := ids[2]
	switch taskType {
	case "task":
		_, _, err = ironClient.Tasks.CancelTask(taskID)
	case "schedule":
		_, _, err = ironClient.Schedules.CancelSchedule(taskID)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = ironClient.Codes.DeleteCode(codeID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags

}

func resourceFunctionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ironClient, _, _, err := newIronClient(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("docker_credentials") {
		if _, ok := d.GetOk("docker_credentials"); ok {
			if _, err := dockerLogin(ironClient, d); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}

func resourceFunctionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	ironClient, ironConfig, _, err := newIronClient(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceFunctionRead.newIronClient: %w", err))
	}
	ids := strings.Split(d.Id(), "-")
	taskType := ids[0]
	taskID := ids[1]
	codeID := ids[2]

	code, _, err := ironClient.Codes.GetCode(codeID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceFunctionRead.GetCode: %w", err))
	}
	switch taskType {
	case "task":
		task, _, err := ironClient.Tasks.GetTask(taskID)
		if err != nil {
			return diag.FromErr(err)
		}
		_, _ = config.Debug("ProjectID: %v\nType: %v, ID: %v\nCode: %v\n", ironConfig.ProjectID, task, code)
	case "schedule":
		schedule, _, err := ironClient.Schedules.GetSchedule(taskID)
		if err != nil {
			return diag.FromErr(err)
		}
		_, _ = config.Debug("ProjectID: %v\nType: %v, ID: %v\nCode: %v\n", ironConfig.ProjectID, schedule, code)

	}
	return diags
}

type payload struct {
	Version  string            `json:"version"`
	Type     string            `json:"type"`
	Token    string            `json:"token,omitempty"`
	Upstream string            `json:"upstream,omitempty"`
	Env      map[string]string `json:"env,omitempty"`
	Cmd      []string          `json:"cmd,omitempty"`
}

func resourceFunctionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ironClient, ironConfig, modConfig, err := newIronClient(d)
	if err != nil {
		return diag.FromErr(err)
	}
	taskType := "schedule"
	schedule, isSchedule, err := getSchedule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isSchedule {
		taskType = "task"
	}
	name := d.Get("name").(string)
	dockerImage := d.Get("docker_image").(string)
	if _, ok := d.GetOk("docker_credentials"); ok {
		if _, err := dockerLogin(ironClient, d); err != nil {
			return diag.FromErr(err)
		}
	}
	if ironConfig == nil || len(ironConfig.ClusterInfo) == 0 {
		return diag.FromErr(fmt.Errorf("invalid iron config: %v", ironConfig))
	}
	codeName := fmt.Sprintf("tf-%s", name)
	createdCode, _, err := ironClient.Codes.CreateOrUpdateCode(iron.Code{
		Name:      codeName,
		Image:     dockerImage,
		ProjectID: ironConfig.ProjectID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	encryptedPayload, err := preparePayload(taskType, *modConfig, d, ironConfig)

	if err != nil {
		_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
		return diag.FromErr(err)
	}
	switch taskType {
	case "task":
		task, _, err := ironClient.Tasks.QueueTask(iron.Task{
			CodeName: codeName,
			Payload:  encryptedPayload,
			Cluster:  ironConfig.ClusterInfo[0].ClusterID,
		})
		if err != nil {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s-%s-%s", taskType, task.ID, createdCode.ID))
	case "schedule":
		schedule.CodeName = codeName
		schedule.Payload = encryptedPayload
		schedule.Cluster = ironConfig.ClusterInfo[0].ClusterID
		schedule, _, err = ironClient.Schedules.CreateSchedule(*schedule)
		if err != nil {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s-%s-%s", taskType, schedule.ID, createdCode.ID))
	}
	return diags
}

func preparePayload(taskType string, modConfig map[string]interface{}, d *schema.ResourceData, config *iron.Config) (string, error) {
	command := []string{"/app/server"}
	if list, ok := d.Get("command").([]interface{}); ok && len(list) > 0 {
		command = []string{}
		for i := 0; i < len(list); i++ {
			command = append(command, list[i].(string))
		}
	}
	environment := getEnvironment(d)

	payload := payload{
		Version:  "1",
		Type:     taskType,
		Token:    modConfig["siderite_token"].(string),
		Upstream: modConfig["siderite_upstream"].(string),
		Cmd:      command,
		Env:      environment,
	}
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return "", fmt.Errorf("preparePayload: %w", err)
	}
	return iron.EncryptPayload([]byte(config.ClusterInfo[0].Pubkey), payloadJSON)
}

func getEnvironment(d *schema.ResourceData) map[string]string {
	environment := make(map[string]string)
	if e, ok := d.GetOk("environment"); ok {
		env, ok := e.(map[string]interface{})
		if !ok {
			return map[string]string{}
		}
		for k, v := range env {
			environment[k] = v.(string)
		}
	}
	return environment
}

func dockerLogin(ironClient *iron.Client, d *schema.ResourceData) (bool, error) {
	name := d.Get("name").(string)
	dockerImage := d.Get("docker_image").(string)
	ref, err := reference.ParseNormalizedNamed(dockerImage)
	if err != nil {
		return false, fmt.Errorf("error normalizing docker [%s]: %w", dockerImage, err)
	}
	registry := ""
	if str := strings.Split(ref.Name(), "/"); len(str) > 1 {
		registry = str[0]
	}
	v, ok := d.GetOk("docker_credentials")
	if !ok { // No credentials after all
		return true, nil
	}
	vv := v.(map[string]interface{})
	username := vv["username"].(string)
	password := vv["password"].(string)
	ok, _, err = ironClient.Codes.DockerLogin(iron.DockerCredentials{
		Email:         fmt.Sprintf("terraform-%s@localhost.localdomain", name),
		Username:      username,
		Password:      password,
		ServerAddress: registry,
	})
	if !ok {
		return false, fmt.Errorf("invalid docker credentials: %w", err)
	}
	return true, nil
}

func getSchedule(d *schema.ResourceData) (*iron.Schedule, bool, error) {
	schedule, ok := d.Get("schedule").([]interface{})
	if !ok {
		return nil, false, nil
	}
	if len(schedule) == 0 {
		return nil, false, nil
	}
	scheduleResource, ok := schedule[0].(map[string]interface{})
	if !ok {
		return nil, false, nil
	}
	startAt, err := time.Parse(time.RFC3339, scheduleResource["start"].(string))
	if err != nil {
		return nil, false, err
	}
	runEvery, err := calcRunEvery(scheduleResource["run_every"].(string))
	if err != nil {
		return nil, false, err
	}
	ironSchedule := iron.Schedule{
		StartAt:  &startAt,
		RunEvery: runEvery,
	}
	return &ironSchedule, true, nil
}

func calcRunEvery(runEvery string) (int, error) {
	var unit string
	var value int
	scanned, err := fmt.Sscanf(runEvery, "%d%s", &value, &unit)
	if err != nil {
		return 0, fmt.Errorf("runEvery scan [%s]: %w", runEvery, err)
	}
	if scanned != 2 {
		return 0, fmt.Errorf("invalid run_every format: %s", runEvery)
	}
	seconds := 0
	switch unit {
	case "s":
		seconds = value
	case "m":
		seconds = 60 * value
	case "h":
		seconds = 3600 * value
	case "d":
		seconds = 86400 * value
	default:
		return 0, fmt.Errorf("unit '%s' not supported", unit)
	}
	if seconds < 60 {
		return 0, fmt.Errorf("a value less than 60 seconds is not supported")
	}
	return seconds, nil
}

func newIronClient(d *schema.ResourceData) (*iron.Client, *iron.Config, *map[string]interface{}, error) {
	backend, ok := d.Get("backend").([]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("expected array of 'backend' config")
	}
	backendResource, ok := backend[0].(map[string]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("unexpected backend format")
	}
	backendType, ok := backendResource["type"].(string)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid backend type")
	}
	if backendType != "iron" {
		return nil, nil, nil, fmt.Errorf("expected backed type of 'iron'")
	}
	config, ok := backendResource["credentials"].(map[string]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid or missing iron config credentials")
	}
	/*
		var config map[string]string
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return nil, nil, nil, fmt.Errorf("error parsing iron config: %w", err)
		}
	*/
	ironConfig := iron.Config{
		Email:     config["email"].(string),
		Password:  config["password"].(string),
		Project:   config["project"].(string),
		ProjectID: config["project_id"].(string),
		Token:     config["token"].(string),
		UserID:    config["user_id"].(string),
		ClusterInfo: []iron.ClusterInfo{
			{
				ClusterID:   config["cluster_info_0_cluster_id"].(string),
				ClusterName: config["cluster_info_0_cluster_name"].(string),
				Pubkey:      config["cluster_info_0_pubkey"].(string),
				UserID:      config["cluster_info_0_user_id"].(string),
			},
		},
	}

	client, err := iron.NewClient(&ironConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("iron.NewClient: %w", err)
	}
	return client, &ironConfig, &config, nil
}
