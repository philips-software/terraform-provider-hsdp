---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_service_agent

Retrieve details of a ServiceAgent

## Example Usage

```hcl
data "hsdp_connect_mdm_service_agent" "postgres_service_agent" {
   name = "postgreserviceagent"
}
```

```hcl
output "service_agent_configuration" {
   value = data.hsdp_connect_service_agents.postgres_service_agent.configuration
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service agent

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceAgent ID
* `description` - The ServiceAgent description
* `configuration` - The service agent configuration
* `data_subscriber_id` - The service agent data subscriber ID
* `api_version_supported` - The supported api versions of the service agent
* `data_subscriber_id` - The data subscriber ID
* `authentication_method_ids` - The list of authentication methods
