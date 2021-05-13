# hsdp_iam_service
Provides details of a given HSDP IAM service.

## Example Usage

```hcl
data "hsdp_iam_service" "my_service" {
  service_id = "service1@myorg.philips-healthsuite.com"
}
```

```hcl
output "service1_uuid" {
   value = data.hsdp_iam_service.my_service.uuid
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The service ID of the service in HSDP IAM

## Attributes Reference

The following attributes are exported:

* `uuid` - The UUID of the user
* `organization_id` - The organization ID associated with this service
* `application_id` - The application ID associated with this sefvice
* `name` - The name of the service
* `description` - The service description
* `expires_on` - When the service expires (string)
* `scopes` - The scopes assigned to this service
* `default_scopes` - The default scopes of this service

## Error conditions

If the service does not fall under the given organization administration lookup may fail. In that case the lookup will return the following error

`responseCode: 4010`
