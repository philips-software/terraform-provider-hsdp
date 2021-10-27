# hsdp_docker_repository

Retrieves information on a HSDP Docker repository

~> This resource only works when `HSDP_UAA_USERNAME` and `HSDP_UAA_PASSWORD` values matching provider arguments are set

## Example usage

```hcl
data "hsdp_docker_namespace" "apps" {
  name = "apps"
}

data "hsdp_docker_repository" "caddy" {
  namespace_id = data.hsdp_docker_namespace.apps.id
  name         = "caddy"
}

output "tags" {
  value = data.hsdp_docker_repository.caddy.tags
}
```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Required) The organization users should belong to
* `name` - (Required) Filter users on verified email state

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `short_description` - A short description of the repository
* `full_description` - A longer description, supporting markdown
* `ids` - The ids of the tags
* `tags` - The list of tags names
* `updated_at` - The update timestamp
* `compressed_sizes` - The compressed size
* `image_digests` - The SHA256 image digest
* `image_ids` - The image ids
* `num_pulls` - The pulls per tag
* `total_pulls` - The total number of pulls for this repo
