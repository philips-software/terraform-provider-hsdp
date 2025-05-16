---
subcategory: "AI Workspace"
---

# hsdp_ai_workspace_service_instance

-> **Deprecation Notice** This data source is deprecated and will be removed in an upcoming release of the provider

Retrieve details of an existing HSDP AI Workspace service instance.

## Example Usage

```hcl
data "hsdp_config" "workspace" {
  service = "workspace"
}

data "hsdp_ai_workspace_service_instance" "workspace" {
  base_url        = data.hsdp_config.workspace.url
  organization_id = var.workspace_tenant_org_id
}
```

## Argument Reference

The following arguments are supported:

* `base_url` - (Required) the base URL of the Workspace deployment. This can be auto-discovered and/or provided by HSDP.
* `organization_id` - (Required) the Tenant IAM organization.

## Attributes Reference

The following attributes are exported:

* `endpoint` - The Workspace endpoint URL
