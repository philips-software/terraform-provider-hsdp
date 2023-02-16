---
subcategory: "AI Workspace"
page_title: "HSDP: hsdp_ai_workspace_compute_target"
description: |-
  Manages HSDP AI Workspaces Compute targets
---

# hsdp_ai_compute_target

Manages HSDP AI Workspace compute targets

## Example usage

```hcl
data "hsdp_config" "workspace" {
  service = "workspace"
}

data "hsdp_ai_workspace_service_instance" "ws" {
  base_url        = data.hsdp_config.workspace.url
  organization_id = var.workspace_tenant_org_id
}

resource "hsdp_ai_inference_compute_target" "target" {
  endpoint = data.hsdp_ai_inference_service_instance.inference.endpoint

  name          = "test-target"
  description   = "First Compute Target"
  instance_type = "ml.p3.16xlarge"
  storage       = 20

  depends_on = [module.ai-inference-onboarding]
}

resource "hsdp_ai_workspace" "workspace1" {
  endpoint = data.hsdp_ai_inference_service_instance.inference.endpoint
  
  name          = "Workspace 1"
  description   = "Test workspace for my world changing algorithm"
 
  compute_target  {
    identifier = hsdp_ai_workspace_compute_target.target.id
  }
  
  source_code  {
    url       = "git@github.com:loafoe/algo.git"
    branch    = "main"
    commit_id = "e1f9366"
    ssh_key   = "..."
  }
  
  labels = ["CNN"]
}
```

## Argument reference

The following arguments are supported:

* `endpoint` - (Required) The AI Inference instance endpoint
* `name` - (Required) The name of the Model
* `instance_type` - (Required) The instance type
* `storage` - (Required) The storage to allocate

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the Model
* `created` - The date this Model  was created
* `created_by` - Who created the Model
