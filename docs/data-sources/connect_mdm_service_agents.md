---
subcategory: "Master Data Management (MDM)"
---

# hsdp_connect_mdm_service_agents

Retrieve details of ServiceAgents

## Example Usage

```hcl
data "hsdp_connect_mdm_service_agents" "all" {
}
```

```hcl
output "service_agent_names" {
   value = data.hsdp_connect_service_agents.all.names
}
```

## Attributes Reference

The following attributes are exported:

* `ids` - The ServiceAgent IDs
* `names` - The ServiceAgent names
* `descriptions` - The ServiceAgent descriptions
* `configurations` - The service agent configurations
* `data_subscriber_ids` - The service agent data subscriber IDs
* `supported_api_versions` - The supported api versions of the service agents
