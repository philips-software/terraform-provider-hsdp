# hsdp_docker_repository

Manages a HSDP Docker repository

~> This resource only works when `HSDP_UAA_USERNAME` and `HSDP_UAA_PASSWORD` values matching provider arguments are set

## Example usage

```hcl
resource "hsdp_docker_namespace" "apps" {
  name = "apps"
}

resource "hsdp_docker_repository" "caddy" {
  namespace_id = hsdp_docker_namespace.apps.id
  name         = "caddy"
  
  short_description = "Caddy server image" 
  full_description  = "A copy of the official Caddy Docker image"
}
```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Required) The organization users should belong to
* `name` - (Required) Filter users on verified email state
* `short_description` - (Optional) A short description of the repository (100 chars max)
* `full_description` - (Optional) A longer description, supporting markdown

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ids of the repository
* `tags` - The list of tag names
* `updated_at` - The update timestamp
* `compressed_sizes` - The compressed size of the tags
* `image_digests` - The SHA256 image digest of the tags
* `image_ids` - The image ids of the tags
* `num_pulls` - The pulls per tag
* `total_pulls` - The total number of pulls for this repo
