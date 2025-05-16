---
subcategory: "AI Workspace"
page_title: "HSDP: hsdp_ai_workspace"
description: |-
  Manages HSDP AI Workspaces
---

# hsdp_ai_workspace

-> **Deprecation Notice** This resource is deprecated and will be removed in an upcoming release of the provider

Manages HSDP AI Workspace instances

## Example usage

```hcl
data "hsdp_config" "workspace" {
  service = "workspace"
}

data "hsdp_ai_workspace_service_instance" "ws" {
  base_url        = data.hsdp_config.workspace.url
  organization_id = var.workspace_tenant_org_id
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
* `compute_target` - (Required)
  * `reference` - (Required) Reference to the compute target
* `source_code` - (Required)
* `labels` - (Optional)
* `artifact_path` - (Optional)
* `description` - (Optional) Description of the Compute Target
* `additional_configuration` - (Optional)

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the Model
* `created` - The date this Model  was created
* `created_by` - Who created the Model
