---
subcategory: "Notification"
page_title: "HSDP: hsdp_notification_topic"
description: |-
  Manages HSDP Notifation Topic resources
---

# hsdp_notification_topic

Create and manage HSDP Notification Topic resources

## Example usage

```hcl
resource "hsdp_notification_topic" "topic" {
  name =  "Topic1"
  producer_id =  "036c8a21-6906-4485-b2e7-e31883d8f9ed"
  scope =  "public"
  allowed_scopes = [
    "*.*.*.NotificationTest"
  ]
  is_auditable =  true
  description = "topic description"
}
```

## Argument reference

* `name` - (Required) The name of a topic. The topic name length is restricted to a maximum of 256 characters. The special characters allowed are `-` and `_`.
* `producer_id` - (Required) The UUID of the producer
* `scope` - (Required) The intended scope of this topic. Can be either `public` or `private`
* `allowed_scopes` - (Required, list(string)) Validates whether the subscriber can access the topic

  One topic can have multiple allowedScopes, depending on the number of subscribers. The current release only validates at the organization level and the scope attached to the service account/client.

  The following pattern is valid for the allowedScopes field:
`*.*.*.*` -> `ORGANIZATION.PROPOSITION.APPLICATION.SCOPE`

  Expected values for public topics: `*.*.*.somescopevalue`. Public topics will not allow ? for organization value or * in scopevalues.

  Expected values for private topics:

  `?.*.*.*` - Both producer and subscriber should belong to the same organization. Scope need not be present or it can be of any value.

  `?.*.*.somescope` - Both producer and subscriber should belong to the same organization. Subscriber should have the same scope as mentioned in the allowedScopes property.

  `Org1.*.*.somescope` - Subscriber should belong to Org1 and should have the scope as specified.
Private topics will not allow * for the organization value.

* `is_auditable` - (Optional) Indicates whether the topic has to be audited whenever messages are published to it. Default value is `false`. User has to set to `true` for audit to happen.
* `description` - (Optional) The intended usage of this topic
* `principal` - (Optional) The optional principal to use for this resource
  * `service_id` - (Optional) The IAM service ID
  * `service_private_key` - (Optional) The IAM service private key to use
  * `region` - (Optional) Region to use. When not set, the provider config is used
  * `environment` - (Optional) Environment to use. When not set, the provider config is used
  * `endpoint` - (Optional) The endpoint URL to use if applicable. When not set, the provider config is used

## Attribute reference

* `id` - The topic ID

## Importing

Importing a HSDP Notification topic is supported.
