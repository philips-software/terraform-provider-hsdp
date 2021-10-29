---
subcategory: "S3 Credentials"
---

# hsdp_s3creds_access

Gets credentials for an S3 Credentials access

## Example Usage

```hcl
data "hsdp_s3creds_access" "my_access" {
   product_key = var.product_key
   username = "my_iam_login"
   password = "MyP@ssw0rd"
}
```

```hcl
output "s3_credentials" {
   value = data.hsdp_s3creds_access.my_access.access
}
```

## Attributes Reference

The following attributes are exported:

* `access` - JSON response to access request

Example output:

```json
[
  {
    "allowed": {
      "resources": [
        "978abfcc-6327-4373-86b4-3eb4ec8cce0f/*"
      ],
      "actions": [
        "GET",
        "PUT",
        "LIST",
        "DELETE",
        "ALL_OBJECT"
      ]
    },
    "credentials": {
      "accessKey": "PV86FAKEdquKdxDeTZ4s",
      "secretKey": "6qqXSECRETZqlP6fhkAuiIAdQyv2pvwL5mAQyOpc",
      "sessionToken": "AnotherFakeSecretTokenThPRGFvODdRVUdjWGl4bzM3WERqQnZ1bDMya3JxdlNpb3FYM01MSFdDZkZRSmZ4VGZoa05qUEJrNUdzRUZ0U3BKRHo1T0g0OGZ3bXd3bWowUzVxYTFLaG9LcnB2YWxHUXZuUTdjTks5VXNQMWVVbXhyQWhvcERidzRodkxMSWh2S25CTFgwZFBTU2ppUkc1ZlJHRGlhRHo3dnNrUHFFZnFYd0o3ZWFZYTlIM1ZMUk9CZ0JjdmgzaHIyU1lOTkNmWmdudEhtN1k2eGw0dWlFVWpNY3dTVjFRQklnamlzbGNPQmJkSmgyb0IyVlZzOU9NaEtHdmFNMFVqa0M5OEs4OWxCWnY4cEs4ZGtUZVhIVngzRjB5eUtJRVFqYVlxVE9PZjNIVXREYUJtMlh2Wk1CeW1zZXgza3RZVFhpRXBKeDBNMGpWbnVwN1NQbnpGVnFQYkJOdlNaZFZZcjFRR0g5Z2V5MENEcGFPNjJyTXVrbXZOd2E5ekFvM0NuUGVZMHdtQmk5a09NcE84ZTlsb21COVIwOER6WlNVdWhPUTdEanZGcFdlTXZ5TkUzajBBZERPSUJ5c1RrYVhZUUVoMUNTWWV4SVBSSXNjVHJOS1lqVmNLTXg0N3NMN0dnU2R4eGtrQlFoNU1vVG1FeXpQcTdFWGVpVXgyM1cK",
      "expires": "2019-02-20T20:58:10.000",
      "bucket": "cf-s3-eb78633b-7833-4953-aa58-cee7d854812b"
    }
  }
]
```
