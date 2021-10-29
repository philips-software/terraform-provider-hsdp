---
subcategory: "AI Workspace"
---

# hsdp_ai_workspace_compute_targets

Retrieves AI Workspace Compute Targets

## Example usage

```hcl
data "hsdp_config" "workspace" {
  service = "workspace"
}

data "hsdp_ai_workspace_instance" "workspace" {
  base_url        = data.hsdp_config.workspace.url
  organization_id = var.workspace_tenant_org_id
}

data "hsdp_ai_workspace_compute_targets" "targets" {
  endpoint = data.hsdp_ai_workspace_instance.workspace.endpoint
}
```

## Argument reference

* `endpoint`- (Required) the AI Workspace endpoint

## Attribute reference

The following attributes are exported:

* `ids` -  The list of container host IDs
* `names` - The list of container host names. This matches up with the `ids` list index.
