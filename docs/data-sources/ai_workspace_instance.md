# hsdp_ai_workspace_instance

Retrieve details of an existing HSDP AI Workspace instance.

## Example Usage

```hcl
data "hsdp_config" "workspace" {
  service = "workspace"
}

data "hsdp_ai_workspace_instance" "workspace" {
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
