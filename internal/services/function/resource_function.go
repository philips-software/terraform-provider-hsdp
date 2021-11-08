package function

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
	"github.com/philips-labs/ferrite/server"
	siderite "github.com/philips-labs/siderite/models"
	"github.com/philips-software/go-hsdp-api/iron"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/tools"
	"github.com/robfig/cron/v3"
)

const (
	aLongTime = 86400 * 365 * 30
)

func ResourceFunction() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 7,
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
			},
			"docker_credentials": {
				Type:      schema.TypeMap,
				Optional:  true,
				Sensitive: true,
			},
			"command": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"environment": {
				Type:      schema.TypeMap,
				Optional:  true,
				Sensitive: true,
			},
			"run_every": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"schedule"},
			},
			"start_at": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"schedule"},
			},
			"schedule": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"run_every", "start_at"},
				ValidateDiagFunc: tools.ValidateCron,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1800,
			},
			"backend": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"sync_endpoint": {
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
	if len(ids) < 2 {
		d.SetId("") // Malformed
		return diags
	}
	codeID := ids[0]
	_ = ids[1]                                      // signature
	_, _, err = ironClient.Codes.DeleteCode(codeID) // Deleting a code cascade deletes schedules as well, jobs done!
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags

}

func resourceFunctionUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ironClient, ironConfig, modConfig, err := newIronClient(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	c := m.(*config.Config)

	if d.HasChange("docker_credentials") {
		if _, ok := d.GetOk("docker_credentials"); ok {
			if _, err := dockerLogin(ironClient, d); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// ID Format: {codeID}-{signature}
	ids := strings.Split(d.Id(), "-")
	if len(ids) < 2 {
		d.SetId("") // Malformed
		return diags
	}
	codeID := ids[0]
	signature := ids[1]
	name := d.Get("name").(string)
	codeName := fmt.Sprintf("%s-%s", name, signature)
	if d.HasChange("docker_image") {
		code, _, err := ironClient.Codes.GetCode(codeID)
		if err != nil {
			d.SetId("")
			return diags
		}
		code.Image = d.Get("docker_image").(string)
		_, resp, err := ironClient.Codes.CreateOrUpdateCode(*code)
		if err != nil {
			return diag.FromErr(fmt.Errorf("CreateOrUpdateCode(%v): %w", code, err))
		}
		if resp.StatusCode != http.StatusOK {
			return diag.FromErr(fmt.Errorf("failed to update code: %d", resp.StatusCode))
		}
	}

	if d.HasChange("schedule") || d.HasChange("command") ||
		d.HasChange("run_every") || d.HasChange("environment") ||
		d.HasChange("start_at") || d.HasChange("docker_image") || d.HasChange("timeout") {
		schedules, _, err := ironClient.Schedules.GetSchedulesWithCode(codeName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("GetSchedulesWithCode(%s): %w", codeName, err))
		}
		// Create new schedules
		diags = createSchedules(ironClient, ironConfig, *modConfig, d, codeName, codeID, signature)
		if len(diags) > 0 {
			return diags
		}
		// Clear old ones
		for _, s := range *schedules {
			_, _, _ = ironClient.Schedules.CancelSchedule(s.ID)
		}
	}
	_, _ = c.Debug("ProjectID: %v\nSignature: %v\nCode: %v\n", ironConfig.ProjectID, signature, codeID)
	return diags
}

func resourceFunctionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)
	ironClient, ironConfig, _, err := newIronClient(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceFunctionRead.newIronClient: %w", err))
	}
	// ID Format: {codeID}-{signature}
	ids := strings.Split(d.Id(), "-")
	codeID := ids[0]
	signature := ids[1]

	code, _, err := ironClient.Codes.GetCode(codeID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceFunctionRead.GetCode: %w", err))
	}
	if code == nil || code.ID != codeID {
		_, _ = c.Debug("could not find code with ID: %s. marking resource as gone\n", codeID)
		_ = d.Set("docker_image", "")
	} else {
		_ = d.Set("docker_image", code.Image)
	}
	schedules, _, err := ironClient.Schedules.GetSchedulesWithCode(code.Name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceFunctionRead.GetSchedulesWithCode: %w", err))
	}

	// Check schedules
	if len(*schedules) == 0 {
		_, _ = c.Debug("no schedules found for code '%s'. marking resource as gone\n", code.Name)
		_ = d.Set("docker_image", "")
		return diags
	}
	_, _ = c.Debug("ProjectID: %v\nSignature: %v\nCode: %v\nSchedules: %d\n", ironConfig.ProjectID, signature, codeID, len(*schedules))
	return diags
}

func resourceFunctionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ironClient, ironConfig, modConfig, err := newIronClient(d, m)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("invalid iron discovery: %v", ironConfig))
	}
	codeName := fmt.Sprintf("%s-%s", name, signature)
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
	diags := createSchedules(ironClient, ironConfig, *modConfig, d, codeName, createdCode.ID, signature)
	if len(diags) > 0 {
		return diags
	}

	_ = d.Set("token", (*modConfig)["siderite_token"])
	_ = d.Set("auth_type", (*modConfig)["siderite_auth_type"])
	return resourceFunctionRead(ctx, d, m)
}

func createSchedules(ironClient *iron.Client, ironConfig *iron.Config, modConfig map[string]string, d *schema.ResourceData, codeName, codeID, signature string) diag.Diagnostics {
	var diags diag.Diagnostics

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
	encryptedSyncPayload, encryptedAsyncPayload, err := preparePayloads(taskType, modConfig, d, ironConfig)
	if err != nil {
		_, _, _ = ironClient.Codes.DeleteCode(codeID)
		return diag.FromErr(err)
	}
	timeout := d.Get("timeout").(int)
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
		_, resp, err := ironClient.Schedules.CreateSchedule(cronSchedule)
		if err != nil || resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(codeID)
			if err == nil {
				return diag.FromErr(fmt.Errorf("create CRON schedule failed with code %d", resp.StatusCode))
			}
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s-%s", codeID, signature))
	case "function":
		cfg := siderite.CronPayload{
			EncryptedPayload: encryptedSyncPayload,
			Type:             "sync",
		}
		jsonPayload, _ := json.Marshal(cfg)
		syncSchedule = &iron.Schedule{
			CodeName: codeName,
			Payload:  string(jsonPayload),
			Cluster:  ironConfig.ClusterInfo[0].ClusterID,
			StartAt:  &startAt,
			RunEvery: aLongTime,
			Timeout:  timeout,
		}
		_, resp, err := ironClient.Schedules.CreateSchedule(*syncSchedule)
		if err != nil {
			_, _, _ = ironClient.Codes.DeleteCode(codeID)
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(codeID)
			return diag.FromErr(fmt.Errorf("creating sync schedule failed with code %d", resp.StatusCode))
		}
		cfg = siderite.CronPayload{
			EncryptedPayload: encryptedAsyncPayload,
			Type:             "async",
		}
		jsonPayload, _ = json.Marshal(cfg)
		asyncSchedule = &iron.Schedule{
			CodeName: codeName,
			Payload:  string(jsonPayload),
			Cluster:  ironConfig.ClusterInfo[0].ClusterID,
			StartAt:  &startAt,
			RunEvery: aLongTime,
			Timeout:  timeout,
		}
		_, resp, err = ironClient.Schedules.CreateSchedule(*asyncSchedule)
		if err != nil {
			_, _, _ = ironClient.Codes.DeleteCode(codeID)
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(codeID)
			return diag.FromErr(fmt.Errorf("creating async schedule failed with code %d", resp.StatusCode))
		}
		d.SetId(fmt.Sprintf("%s-%s", codeID, signature))
	case "schedule":
		if schedule == nil {
			return diag.FromErr(fmt.Errorf("expected schedule to not be nil"))
		}
		schedule.Iron.CodeName = codeName
		schedule.Iron.Payload = encryptedSyncPayload
		schedule.Iron.Cluster = ironConfig.ClusterInfo[0].ClusterID
		_, resp, err := ironClient.Schedules.CreateSchedule(*schedule.Iron)
		if err != nil || resp.StatusCode != http.StatusOK {
			_, _, _ = ironClient.Codes.DeleteCode(codeID)
			if err == nil {
				return diag.FromErr(fmt.Errorf("create schedule failed with code %d", resp.StatusCode))
			}
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s-%s", codeID, signature))
	}
	if syncSchedule != nil {
		syncEndpoint := fmt.Sprintf("https://%s/function/%s", (modConfig)["siderite_upstream"], codeID)
		_ = d.Set("endpoint", syncEndpoint)
		_ = d.Set("sync_endpoint", syncEndpoint)
	}
	if asyncSchedule != nil {
		asyncEndpoint := fmt.Sprintf("https://%s/async-function/%s", (modConfig)["siderite_upstream"], codeID)
		_ = d.Set("async_endpoint", asyncEndpoint)
	}
	return diags
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

type taskSchedule struct {
	Timeout int
	Iron    *iron.Schedule
	CRON    *string
}

func getSchedule(d *schema.ResourceData) (*taskSchedule, bool, error) {
	timeout := d.Get("timeout").(int)
	// Check for cron
	cronSchedule := d.Get("schedule").(string)
	if cronSchedule != "" {
		_, err := cron.ParseStandard(cronSchedule)
		if err != nil {
			return nil, false, fmt.Errorf("parsing cron field: %w", err)
		}
		return &taskSchedule{CRON: &cronSchedule, Timeout: timeout}, true, nil
	}
	runEverySchedule := d.Get("run_every").(string)
	startAtSchedule := d.Get("start_at").(string)
	if runEverySchedule != "" {
		runEvery, firstRun, err := calcRunEvery(runEverySchedule, startAtSchedule)
		if err != nil {
			return nil, false, err
		}
		ironSchedule := iron.Schedule{
			StartAt:  firstRun,
			RunEvery: runEvery,
			Timeout:  timeout,
		}
		return &taskSchedule{Iron: &ironSchedule}, true, nil
	}
	return nil, false, nil
}

func calcRunEvery(runEvery, startAt string) (int, *time.Time, error) {
	var unit string
	var value int
	scanned, err := fmt.Sscanf(runEvery, "%d%s", &value, &unit)
	if err != nil {
		return 0, nil, fmt.Errorf("runEvery scan [%s]: %w", runEvery, err)
	}
	if scanned != 2 {
		return 0, nil, fmt.Errorf("invalid run_every format: %s", runEvery)
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
		return 0, nil, fmt.Errorf("unit '%s' not supported", unit)
	}
	if seconds < 60 {
		return 0, nil, fmt.Errorf("a value less than 60 seconds is not supported")
	}
	now := time.Now()
	firstRun := now
	if startAt != "" {
		userFirstRun, err := time.Parse(time.RFC3339, startAt)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid start_at time: %s", startAt)
		}
		firstRun = userFirstRun
	}
	if firstRun.Before(now) { // In the past so figure out the firstRun time
		timeComponent := firstRun.Sub(firstRun.Truncate(24 * time.Hour))
		firstRun = now.Truncate(24 * time.Hour).Add(timeComponent)
		if now.After(firstRun) { // Start tomorrow
			firstRun = firstRun.Add(24 * time.Hour)
		}
	}
	return seconds, &firstRun, nil
}

func newIronClient(d *schema.ResourceData, m interface{}) (*iron.Client, *iron.Config, *map[string]string, error) {
	c := m.(*config.Config)
	backend, ok := d.Get("backend").([]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("expected array of 'backend' discovery")
	}
	backendResource, ok := backend[0].(map[string]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("unexpected backend format")
	}

	configMap, ok := backendResource["credentials"].(map[string]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid or missing iron discovery credentials")
	}
	cfg := make(map[string]string)
	for k, v := range configMap {
		if str, ok := v.(string); ok {
			cfg[k] = str
		}
	}
	backendType := cfg["type"]
	if !(backendType == "siderite" || backendType == "ferrite") {
		return nil, nil, nil, fmt.Errorf("expected backend type of ['siderite' | 'ferrite']")
	}
	// Bootstrap ferrite
	if backendType == "ferrite" {
		token := cfg["token"]
		bootstrap, err := server.Bootstrap(cfg["base_url"], token)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error bootstrapping ferrite: %w", err)
		}
		// Inject bootstrap data
		cfg["project"] = bootstrap.ProjectID
		cfg["project_id"] = bootstrap.ProjectID
		cfg["cluster_info_0_cluster_id"] = bootstrap.ClusterID
		cfg["cluster_info_0_pubkey"] = bootstrap.PublicKey
	}
	ironConfig := iron.Config{
		BaseURL:   cfg["base_url"],
		Email:     cfg["email"],
		Password:  cfg["password"],
		Project:   cfg["project"],
		ProjectID: cfg["project_id"],
		Token:     cfg["token"],
		UserID:    cfg["user_id"],
		DebugLog:  c.DebugLog,
		ClusterInfo: []iron.ClusterInfo{
			{
				ClusterID:   cfg["cluster_info_0_cluster_id"],
				ClusterName: cfg["cluster_info_0_cluster_name"],
				Pubkey:      cfg["cluster_info_0_pubkey"],
				UserID:      cfg["cluster_info_0_user_id"],
			},
		},
	}

	client, err := iron.NewClient(&ironConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("iron.NewClient: %w", err)
	}
	return client, &ironConfig, &cfg, nil
}
