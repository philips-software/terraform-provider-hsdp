---
subcategory: "Blob Repository (BLR)"
page_title: "HSDP: hsdp_blr_blob_store_policy"
description: |-
  Manages HSDP Connect Blob Store Repository Policies
---

# hsdp_blr_blob_store_policy

Create and manage Blob Repository Policies

## Example Usage

```hcl
resource "hsdp_blr_blob_store_policy" "policy" {
  statement {
    effect    = "Allow"
    action    = ["GET", "PUT", "DELETE"]
    principal = ["prn:hsdp:iam:${data.hsdp_iam_org.myorg.id}:${hsdp_connect_mdm_proposition.first.guid}:User/*"]
    resource  = ["${hsdp_blr_bucket.store.name}/*"]
  }
}
```

## Argument Reference

The following arguments are available:

* `statement` - (Required)
    * `effect` - (Required, string) Effect of policy [`Allow`, `Deny`]
    * `action` - (Required, list(string)) Allowed methods: [`GET`, `PUT`, `DELETE`]
    * `principal` - (Required, list(string)) The principals the policy applies to
    * `resource` - (Required, list(string)) The resources the policy applies to

## Attributes reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID reference of the service action (format: `BlobStorePolicy/${GUID}`)
* `guid` - The GUID of the bucket
