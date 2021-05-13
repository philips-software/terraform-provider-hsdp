# hsdp_iam_user
Provides details of a given HSDP IAM user. 

>Typically, this resource is used to only test account. We highly recommend using the IAM Self serviceUI which HSDP provides for day to day user management tasks

## Example Usage

```hcl
data "hsdp_iam_user" "john" {
  username = "john.doe@1e100.io"
}
```

```hcl
output "johns_uuid" {
   value = data.hsdp_iam_user.john.uuid
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The username/email of the user in HSDP IAM

## Attributes Reference

The following attributes are exported:

* `uuid` - The UUID of the user

## Error conditions

If the user does not fall under the given organization administration lookup may fail. In that case the lookup will return the following error

`responseCode: 4010`
