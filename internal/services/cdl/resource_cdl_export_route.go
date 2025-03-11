package cdl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dip-software/go-dip-api/cdl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func exportRouteLabelSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"label_name": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"approval_required": {
					Type:     schema.TypeBool,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func exportRouteDataObjectDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resource_type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"associated_labels": exportRouteLabelSchema(),
			},
		},
	}
}

func sourceResearchStudyDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"source_research_study_endpoint": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"allowed_data_objects": exportRouteDataObjectDetailsSchema(),
			},
		},
	}
}

func serviceAccountDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_id": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"private_key": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"access_token_endpoint": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"token_endpoint": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func ResourceCDLExportRoute() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceCDLExportRouteCreate,
		ReadContext:   resourceCDLExportRouteRead,
		DeleteContext: resourceCDLExportRouteDelete,

		Schema: map[string]*schema.Schema{
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"export_route_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_research_study": sourceResearchStudyDetailsSchema(),
			"auto_export": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"destination_research_study_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_account_details": serviceAccountDetailsSchema(),
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"updated_on": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCDLExportRouteDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*config.Config)

	endpoint := d.Get("cdl_endpoint").(string)
	exportRouteId := d.Id()

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	resp, err := client.ExportRoute.DeleteExportRouteByID(exportRouteId)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}
	return diags
}

func getSourceResearchStudyDetails(d *schema.ResourceData) cdl.ExportResearchStudySource {
	var exportResearchStudySource cdl.ExportResearchStudySource
	if v, ok := d.GetOk("source_research_study"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			sourceResearchStudyField := vi.(map[string]interface{})
			sourceCdlEndpoint := sourceResearchStudyField["source_research_study_endpoint"].(string)
			exportResearchStudySource.Endpoint = sourceCdlEndpoint

			allowedDataObjectsArray := sourceResearchStudyField["allowed_data_objects"].([]interface{})
			var exportDataObjectArray []cdl.ExportDataObject

			for _, i := range allowedDataObjectsArray {
				allowedDataObjectArrayElem := i.(map[string]interface{})
				resourceType := allowedDataObjectArrayElem["resource_type"].(string)
				associatedLabels := allowedDataObjectArrayElem["associated_labels"].([]interface{})

				var exportLabel []cdl.ExportLabel
				for _, l := range associatedLabels {
					labelMap := l.(map[string]interface{})
					labelName := labelMap["label_name"].(string)
					approvalRequired := labelMap["approval_required"].(bool)
					exportLabel = append(exportLabel, cdl.ExportLabel{
						Name:             labelName,
						ApprovalRequired: approvalRequired,
					})
				}
				exportDataObjectArray = append(exportDataObjectArray, cdl.ExportDataObject{
					Type:        resourceType,
					ExportLabel: exportLabel,
				})
			}
			if len(exportDataObjectArray) > 0 {
				exportResearchStudySource.Allowed = &cdl.ExportAllowedField{
					DataObject: exportDataObjectArray,
				}
			}
		}
	}
	return exportResearchStudySource
}

func getServiceAccountDetails(d *schema.ResourceData) cdl.ExportServiceAccount {
	var serviceId, privateKey, accessTokenEndpoint, tokenEndpoint string
	if v, ok := d.GetOk("service_account_details"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			serviceAccountField := vi.(map[string]interface{})
			serviceId = serviceAccountField["service_id"].(string)
			privateKey = serviceAccountField["private_key"].(string)
			accessTokenEndpoint = serviceAccountField["access_token_endpoint"].(string)
			tokenEndpoint = serviceAccountField["token_endpoint"].(string)
		}
	}
	return cdl.ExportServiceAccount{
		CDLServiceAccount: cdl.ExportServiceAccountDetails{
			ServiceID:           serviceId,
			PrivateKey:          privateKey,
			AccessTokenEndPoint: accessTokenEndpoint,
			TokenEndPoint:       tokenEndpoint,
		},
	}
}

func resourceCDLExportRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var exportRouteToCreate cdl.ExportRoute
	endpoint := d.Get("cdl_endpoint").(string)

	exportRouteToCreate.ExportRouteName = d.Get("export_route_name").(string)
	exportRouteToCreate.Description = d.Get("description").(string)
	exportRouteToCreate.DisplayName = d.Get("display_name").(string)

	exportRouteToCreate.Source = cdl.Source{
		CDLResearchStudy: getSourceResearchStudyDetails(d),
	}

	exportRouteToCreate.AutoExport = d.Get("auto_export").(bool)

	exportRouteToCreate.Destination = cdl.Destination{
		CDLResearchStudy: cdl.ExportResearchStudyDestination{
			Endpoint: d.Get("destination_research_study_endpoint").(string),
		},
	}

	exportRouteToCreate.ServiceAccount = getServiceAccountDetails(d)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	createdExportRoute, resp, err := client.ExportRoute.CreateExportRoute(exportRouteToCreate)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode() != http.StatusConflict {
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}
	d.SetId(createdExportRoute.ID)
	return resourceCDLExportRouteRead(ctx, d, m)
}

func resourceCDLExportRouteRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)
	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)

	client, err := c.GetCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()
	id := d.Id()

	exportRoute, _, err := client.ExportRoute.GetExportRouteByID(id)
	if err != nil {
		return diag.FromErr(err)
	} else if exportRoute == nil {
		return diag.FromErr(fmt.Errorf("ExportRoute with ID %s not found", id))
	}

	sourceBytes, err := json.Marshal((*exportRoute).Source)
	if err != nil {
		return diag.FromErr(err)
	}

	destinationBytes, err := json.Marshal((*exportRoute).Destination)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("export_route_name", (*exportRoute).ExportRouteName)
	_ = d.Set("description", (*exportRoute).Description)
	_ = d.Set("display_name", (*exportRoute).DisplayName)
	_ = d.Set("source", string(sourceBytes))
	_ = d.Set("auto_export", (*exportRoute).AutoExport)
	_ = d.Set("destination", string(destinationBytes))
	_ = d.Set("created_by", (*exportRoute).CreatedBy)
	_ = d.Set("created_on", (*exportRoute).CreatedOn)
	_ = d.Set("updated_by", (*exportRoute).UpdatedBy)
	_ = d.Set("updated_on", (*exportRoute).UpdatedOn)

	return diags
}
