# hsdp_config

Retrieve configuration details from services based on region and environment

## Example Usage

```hcl
data "hsdp_config" "iam_us_east_prod" {
  service = "iam"
  region = "us-east"
  environment = "prod"
}
```

```hcl
output "iam_url_us_east_prod" {
   value = data.hsdp_config.iam_us_east_prod.url
}
```
## Argument Reference

The following arguments are supported:

* `service` - (Required) The HSDP service to lookup
* `region` - (Optional) The HSDP region. If not set, defaults to provider level config
* `environment` - (Optional) The HSDP environent. If not set, defaults to provider level config

## Attributes Reference

The following attributes are exported:

* `url` - (string) The (base / API) URL of the service
* `host` - (string) The host of the service
