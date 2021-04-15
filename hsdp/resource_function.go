package hsdp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/siderite"
	"github.com/philips-software/go-hsdp-api/iron"
	"github.com/robfig/cron/v3"
)

const (
	aLongTime = 86400 * 365 * 30
)

func resourceFunction() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 3,
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
							Optional: true,
							Default:  "2020-10-28T00:00:00Z",
							ForceNew: true,
						},
						"run_every": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cron": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							ValidateDiagFunc: validateCron,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1800,
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
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "siderite",
							ForceNew:         true,
							ValidateDiagFunc: validateFunctionBackend,
						},
						"credentials": {
							Type:      schema.TypeMap,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
					},
				},
			},
			"token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"async_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFunctionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ironClient, _, _, err := newIronClient(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	ids := strings.Split(d.Id(), "-")
	schedules := strings.Split(ids[1], ",")
	codeID := ids[2]
	for _, scheduleID := range schedules {
		_, _, err = ironClient.Schedules.CancelSchedule(scheduleID)
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

func resourceFunctionUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ironClient, _, _, err := newIronClient(d, m)
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

func resourceFunctionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*Config)
	ironClient, ironConfig, _, err := newIronClient(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceFunctionRead.newIronClient: %w", err))
	}
	// ID Format: {taskType}-{scheduleID}[,{scheduleID},...]-{codeID}-{signature}
	ids := strings.Split(d.Id(), "-")
	if len(ids) < 4 {
		d.SetId("") // Malformed
		return diags
	}
	taskType := ids[0]
	schedules := strings.Split(ids[1], ",")
	codeID := ids[2]
	_ = ids[3] // signature

	code, _, err := ironClient.Codes.GetCode(codeID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceFunctionRead.GetCode: %w", err))
	}
	if code == nil || code.ID != codeID {
		_, _ = config.Debug("could not find code with ID: %s. marking resource as gone\n", codeID)
		_ = d.Set("docker_image", "")
	} else {
		_ = d.Set("docker_image", code.Image)
	}
	var codeName string
	for _, scheduleID := range schedules {
		schedule, _, err := ironClient.Schedules.GetSchedule(scheduleID)
		if err != nil {
			return diag.FromErr(err)
		}
		if schedule == nil || schedule.ID != scheduleID {
			_, _ = config.Debug("could not find schedule with ID: %s. marking resource as gone\n", scheduleID)
			_ = d.Set("docker_image", "") // Treat missing schedule as destroyed code as well
			return diags
		}
		var scheduleTaskType string
		parts := strings.Split(schedule.CodeName, "-")
		if len(parts) < 2 {
			_, _ = config.Debug("could not match schedule name: '%s'. marking resource as gone\n", schedule.CodeName)
			_ = d.Set("docker_image", "")
			return diags
		}
		scheduleTaskType = parts[0]
		if scheduleTaskType != taskType {
			_, _ = config.Debug("could not match schedule name: '%s'. marking resource as gone\n", schedule.CodeName)
			_ = d.Set("docker_image", "")
			return diags
		}
		codeName = strings.Join(parts[1:len(parts)-1], "-")
	}
	_ = d.Set("name", codeName)
	_, _ = config.Debug("ProjectID: %v\nType: %v, Schedules: %v\nCode: %v\n", ironConfig.ProjectID, taskType, schedules, code)
	return diags
}

func resourceFunctionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ironClient, ironConfig, modConfig, err := newIronClient(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	taskType := "schedule"
	schedule, isSchedule, err := getSchedule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isSchedule {
		taskType = "function"
	}
	if schedule != nil && schedule.CRON != nil {
		taskType = "cron"
	}
	name := d.Get("name").(string)
	dockerImage := d.Get("docker_image").(string)
	if _, ok := d.GetOk("docker_credentials"); ok {
		if _, err := dockerLogin(ironClient, d); err != nil {
			return diag.FromErr(err)
		}
	}
	signature := strings.Replace(uuid.New().String(), "-", "", -1)

	if ironConfig == nil || len(ironConfig.ClusterInfo) == 0 {
		return diag.FromErr(fmt.Errorf("invalid iron config: %v", ironConfig))
	}
	codeName := fmt.Sprintf("%s-%s-%s", taskType, name, signature)
	createdCode, resp, err := ironClient.Codes.CreateOrUpdateCode(iron.Code{
		Name:      codeName,
		Image:     dockerImage,
		ProjectID: ironConfig.ProjectID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return diag.FromErr(fmt.Errorf("code %s not found", dockerImage))
	}
	encryptedSyncPayload, encryptedAsyncPayload, err := preparePayloads(taskType, *modConfig, d, ironConfig)
	if err != nil {
		_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
		return diag.FromErr(err)
	}
	startAt := time.Now().Add(aLongTime * time.Second)
	var syncSchedule *iron.Schedule
	var asyncSchedule *iron.Schedule
	switch taskType {
	case "cron":
		cfg := siderite.CronPayload{
			Schedule:         *schedule.CRON,
			EncryptedPayload: encryptedSyncPayload,
			Timeout:          schedule.Timeout,
		}
		jsonPayload, _ := json.Marshal(cfg)
		cronSchedule := iron.Schedule{
			CodeName: codeName,
			Payload:  string(jsonPayload),
			Cluster:  ironConfig.ClusterInfo[0].ClusterID,
			StartAt:  &startAt,
			RunEvery: aLongTime,
		}
		createdSchedule, resp, err := ironClient.Schedules.CreateSchedule(cronSchedule)
		if err != nil || resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			if err == nil {
				return diag.FromErr(fmt.Errorf("create CRON schedule failed with code %d", resp.StatusCode))
			}
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", taskType, createdSchedule.ID, createdCode.ID, signature))
	case "function":
		syncSchedule = &iron.Schedule{
			CodeName: codeName,
			Payload:  encryptedSyncPayload,
			Cluster:  ironConfig.ClusterInfo[0].ClusterID,
			StartAt:  &startAt,
			RunEvery: aLongTime,
		}
		syncSchedule, resp, err = ironClient.Schedules.CreateSchedule(*syncSchedule)
		if err != nil {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			return diag.FromErr(fmt.Errorf("creating sync schedule failed with code %d", resp.StatusCode))
		}
		asyncSchedule = &iron.Schedule{
			CodeName: codeName,
			Payload:  encryptedAsyncPayload,
			Cluster:  ironConfig.ClusterInfo[0].ClusterID,
			StartAt:  &startAt,
			RunEvery: aLongTime,
		}
		asyncSchedule, resp, err = ironClient.Schedules.CreateSchedule(*asyncSchedule)
		if err != nil {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			_, _, _ = ironClient.Schedules.CancelSchedule(syncSchedule.ID)
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			_, _, _ = ironClient.Schedules.CancelSchedule(syncSchedule.ID)
			return diag.FromErr(fmt.Errorf("creating async schedule failed with code %d", resp.StatusCode))
		}
		d.SetId(fmt.Sprintf("%s-%s,%s-%s-%s", taskType, syncSchedule.ID, asyncSchedule.ID, createdCode.ID, signature))
	case "schedule":
		schedule.Iron.CodeName = codeName
		schedule.Iron.Payload = encryptedSyncPayload
		schedule.Iron.Cluster = ironConfig.ClusterInfo[0].ClusterID
		createdSchedule, resp, err := ironClient.Schedules.CreateSchedule(*schedule.Iron)
		if err != nil || resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(createdCode.ID)
			if err == nil {
				return diag.FromErr(fmt.Errorf("create schedule failed with code %d", resp.StatusCode))
			}
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", taskType, createdSchedule.ID, createdCode.ID, signature))
	}
	if syncSchedule != nil {
		_ = d.Set("endpoint", fmt.Sprintf("https://%s/function/%s", (*modConfig)["siderite_upstream"], syncSchedule.ID))
	}
	if asyncSchedule != nil {
		_ = d.Set("async_endpoint", fmt.Sprintf("https://%s/async-function/%s", (*modConfig)["siderite_upstream"], asyncSchedule.ID))
	}
	_ = d.Set("token", (*modConfig)["siderite_token"])
	_ = d.Set("auth_type", (*modConfig)["siderite_auth_type"])
	return resourceFunctionRead(ctx, d, m)
}

func preparePayloads(taskType string, modConfig map[string]string, d *schema.ResourceData, config *iron.Config) (string, string, error) {
	command := []string{"/app/server"}
	if list, ok := d.Get("command").([]interface{}); ok && len(list) > 0 {
		command = []string{}
		for i := 0; i < len(list); i++ {
			command = append(command, list[i].(string))
		}
	}
	environment := getEnvironment(d)

	payload := siderite.Payload{
		Version:  "1",
		Type:     taskType,
		Token:    modConfig["siderite_token"],
		Upstream: modConfig["siderite_upstream"],
		Auth:     modConfig["siderite_auth_type"],
		Cmd:      command,
		Env:      environment,
		Mode:     "sync",
	}
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return "", "", fmt.Errorf("preparePayload: %w", err)
	}
	syncPayload, err := iron.EncryptPayload([]byte(config.ClusterInfo[0].Pubkey), payloadJSON)
	if err != nil {
		return "", "", fmt.Errorf("preparePayloads.sync: %w", err)
	}
	// Async
	payload.Mode = "async"
	payloadJSON, err = json.Marshal(&payload)
	if err != nil {
		return "", "", fmt.Errorf("preparePayload: %w", err)
	}
	asyncPayload, err := iron.EncryptPayload([]byte(config.ClusterInfo[0].Pubkey), payloadJSON)
	if err != nil {
		return "", "", fmt.Errorf("preparePayloads.async: %w", err)
	}
	return syncPayload, asyncPayload, nil
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

type functionSchedule struct {
	Timeout int
	Iron    *iron.Schedule
	CRON    *string
}

func getSchedule(d *schema.ResourceData) (*functionSchedule, bool, error) {
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
	timeout := scheduleResource["timeout"].(int)
	// Check for cron
	cronSchedule := scheduleResource["cron"].(string)
	if cronSchedule != "" {
		_, err := cron.ParseStandard(cronSchedule)
		if err != nil {
			return nil, false, fmt.Errorf("parsing cron field: %w", err)
		}
		return &functionSchedule{CRON: &cronSchedule, Timeout: timeout}, true, nil
	}
	runEvery, err := calcRunEvery(scheduleResource["run_every"].(string))
	if err != nil {
		return nil, false, err
	}
	ironSchedule := iron.Schedule{
		StartAt:  &startAt,
		RunEvery: runEvery,
		Timeout:  timeout,
	}
	return &functionSchedule{Iron: &ironSchedule}, true, nil
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

func newIronClient(d *schema.ResourceData, m interface{}) (*iron.Client, *iron.Config, *map[string]string, error) {
	c := m.(*Config)

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
	if backendType != "siderite" {
		return nil, nil, nil, fmt.Errorf("expected backend type of 'siderite'")
	}
	configMap, ok := backendResource["credentials"].(map[string]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid or missing iron config credentials")
	}
	config := make(map[string]string)
	for k, v := range configMap {
		if str, ok := v.(string); ok {
			config[k] = str
		}
	}
	ironConfig := iron.Config{
		Email:     config["email"],
		Password:  config["password"],
		Project:   config["project"],
		ProjectID: config["project_id"],
		Token:     config["token"],
		UserID:    config["user_id"],
		DebugLog:  c.DebugLog,
		ClusterInfo: []iron.ClusterInfo{
			{
				ClusterID:   config["cluster_info_0_cluster_id"],
				ClusterName: config["cluster_info_0_cluster_name"],
				Pubkey:      config["cluster_info_0_pubkey"],
				UserID:      config["cluster_info_0_user_id"],
			},
		},
	}

	client, err := iron.NewClient(&ironConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("iron.NewClient: %w", err)
	}
	return client, &ironConfig, &config, nil
}
