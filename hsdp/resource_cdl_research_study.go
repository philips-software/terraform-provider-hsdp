package hsdp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/cdl"
)

func permissionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"user_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"email": {
					Type:     schema.TypeString,
					Required: true,
				},
				"institute_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func resourceCDLResearchStudy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDLResearchStudyCreate,
		ReadContext:   resourceCDLResearchStudyRead,
		UpdateContext: resourceCDLResearchStudyUpdate,
		DeleteContext: resourceCDLResearchStudyDelete,

		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"study_owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ends_at": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressEqualTimeOrMissing,
			},
			"study_manager":  permissionSchema(),
			"data_scientist": permissionSchema(),
			"monitor":        permissionSchema(),
			"uploader":       permissionSchema(),
		},
	}
}

func getPermissionList(field, role string, d *schema.ResourceData) []cdl.RoleRequest {
	var list []cdl.RoleRequest

	if v, ok := d.GetOk(field); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			userID := mVi["user_id"].(string)
			email := mVi["email"].(string)
			instituteID := mVi["institute_id"].(string)
			list = append(list, cdl.RoleRequest{
				IAMUserUUID: userID,
				Email:       email,
				InstituteID: instituteID,
				Role:        role,
			})
		}
	}
	return list
}

func schemaSetToRoleRequest(v interface{}, permission string) []cdl.RoleRequest {
	var list []cdl.RoleRequest
	role := fieldToRole(permission)

	vL := v.(*schema.Set).List()
	for _, vi := range vL {
		mVi := vi.(map[string]interface{})
		userID := mVi["user_id"].(string)
		email := mVi["email"].(string)
		instituteID := mVi["institute_id"].(string)
		list = append(list, cdl.RoleRequest{
			IAMUserUUID: userID,
			Email:       email,
			InstituteID: instituteID,
			Role:        role,
		})
	}
	return list
}

func fieldToRole(field string) string {
	switch field {
	case "monitor":
		return cdl.ROLE_MONITOR
	case "study_manager":
		return cdl.ROLE_STUDYMANAGER
	case "uploader":
		return cdl.ROLE_UPLOADER
	case "data_scientist":
		return cdl.ROLE_DATA_SCIENTIST
	}
	return ""
}

func setPermissionList(field string, studyPermissions cdl.RoleAssignmentResult, d *schema.ResourceData) {
	c := &schema.Set{F: resourceMetricsThresholdHash}
	cL := c.List()
	fRole := fieldToRole(field)
	if fRole == "" {
		return
	}
	// This function just prunes missing items
	for _, role := range studyPermissions {
		for _, vi := range cL {
			mVi := vi.(map[string]interface{})
			if role.IAMUserUUID == mVi["user_id"] {
				for _, r := range role.Roles {
					if r.Role == fRole {
						c.Add(vi)
					}
				}
			}
		}
	}
	_ = d.Set(field, c)
}

func resourceCDLResearchStudyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	title := d.Get("title").(string)
	description := d.Get("description").(string)
	endsAt := d.Get("ends_at").(string)
	studyOwner := d.Get("study_owner").(string)

	createdStudy, resp, err := client.Study.CreateStudy(cdl.Study{
		Title:       title,
		Description: description,
		Period: cdl.Period{
			End: endsAt,
		},
		StudyOwner: studyOwner,
	})
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode != http.StatusConflict {
			return diag.FromErr(err)
		}
		// Search for existing study based on Title
		studies, _, err2 := client.Study.GetStudies(nil)
		if err2 != nil {
			return diag.FromErr(fmt.Errorf("on match attempt during Create conflict: %w", err))
		}
		matchFound := false
		for _, study := range studies {
			if study.Title == title && study.StudyOwner == studyOwner { // TODO: check if this is sufficient for a good match
				d.SetId(study.ID)
				matchFound = true
				break
			}
		}
		if !matchFound {
			return diag.FromErr(err)
		}
	} else {
		d.SetId(createdStudy.ID)
	}

	var perms []cdl.RoleRequest
	for _, f := range []string{"monitor", "uploader", "data_scientist", "study_manager"} {
		perms = append(perms, getPermissionList(f, fieldToRole(f), d)...)
	}
	placeholder := cdl.Study{
		ID: d.Id(),
	}
	for _, r := range perms {
		_, _, _ = client.Study.GrantPermission(placeholder, r)
	}

	return resourceCDLResearchStudyRead(ctx, d, m)
}

func resourceCDLResearchStudyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()

	study, _, err := client.Study.GetStudyByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", study.Title)
	_ = d.Set("description", study.Description)
	_ = d.Set("study_owner", study.StudyOwner)
	_ = d.Set("ends_at", study.Period.End)
	// TODO: we do not read study permissions yet, but really should

	return diags
}

func resourceCDLResearchStudyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()
	study, _, err := client.Study.GetStudyByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("title") || d.HasChange("description") ||
		d.HasChange("ends_at") || d.HasChange("study_owner") {
		study.Title = d.Get("title").(string)
		study.Description = d.Get("description").(string)
		study.Period.End = d.Get("ends_at").(string)
		study.StudyOwner = d.Get("study_owner").(string)

		_, _, err = client.Study.UpdateStudy(*study)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	for _, field := range []string{"data_scientist", "monitor", "uploader", "study_manager"} {
		if d.HasChange(field) {
			o, n := d.GetChange(field)
			vO := schemaSetToRoleRequest(o, field)
			vN := schemaSetToRoleRequest(n, field)
			var toAdd []cdl.RoleRequest
			for _, a := range vN {
				found := false
				for _, b := range vO {
					if a.EqualEnough(b) {
						found = true
					}
				}
				if !found {
					toAdd = append(toAdd, a)
				}
			}
			var toRemove []cdl.RoleRequest
			for _, a := range vO {
				found := false
				for _, b := range vN {
					if a.EqualEnough(b) {
						found = true
					}
				}
				if !found {
					toRemove = append(toRemove, a)
				}
			}
			// $revoke
			for _, r := range toRemove {
				_, resp, err := client.Study.RevokePermission(*study, r)
				if err != nil && resp != nil && resp.StatusCode != http.StatusConflict {
					diags = append(diags, diag.FromErr(err)...)
				}
			}

			// $grant
			for _, r := range toAdd {
				_, resp, err := client.Study.GrantPermission(*study, r)
				if err != nil && resp != nil && resp.StatusCode != http.StatusConflict {
					diags = append(diags, diag.FromErr(err)...)
				}
			}
		}
	}
	return diags
}

func resourceCDLResearchStudyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("") // This is by design currently
	return diags
}
