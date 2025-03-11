package helpers

import (
	"github.com/dip-software/go-dip-api/ai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CollectComputeTarget(d *schema.ResourceData) (ai.ReferenceComputeTarget, diag.Diagnostics) {
	var diags diag.Diagnostics
	var target ai.ReferenceComputeTarget
	if v, ok := d.GetOk("compute_target"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			target = ai.ReferenceComputeTarget{
				Reference:  mVi["reference"].(string),
				Identifier: mVi["identifier"].(string),
			}
		}
	}
	return target, diags
}

func CollectSourceCode(d *schema.ResourceData) (ai.SourceCode, diag.Diagnostics) {
	var diags diag.Diagnostics
	var rce ai.SourceCode
	if v, ok := d.GetOk("source_code"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			rce = ai.SourceCode{
				URL:      mVi["url"].(string),
				Branch:   mVi["branch"].(string),
				CommitID: mVi["commit_id"].(string),
				SSHKey:   mVi["ssh_key"].(string),
			}
		}
	}
	return rce, diags
}

func CollectComputeModel(d *schema.ResourceData) (ai.ReferenceComputeModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model ai.ReferenceComputeModel
	if v, ok := d.GetOk("model"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			model = ai.ReferenceComputeModel{
				Reference:  mVi["reference"].(string),
				Identifier: mVi["identifier"].(string),
			}
		}
	}
	return model, diags
}

func CollectComputeEnvironment(d *schema.ResourceData) (ai.ReferenceComputeEnvironment, diag.Diagnostics) {
	var diags diag.Diagnostics
	var rce ai.ReferenceComputeEnvironment
	if v, ok := d.GetOk("compute_environment"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			rce = ai.ReferenceComputeEnvironment{
				Reference:  mVi["reference"].(string),
				Identifier: mVi["identifier"].(string),
			}
		}
	}
	return rce, diags
}
