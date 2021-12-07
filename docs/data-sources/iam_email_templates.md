---
subcategory: "Master Data Management (MDM)"
---

# hsdp_iam_email_templates

Retrieve details of IAM Email templates configurations for the given IAM Organization

## Example Usage

```hcl
data "hsdp_iam_email_templates" "all" {
  organization_id = data.hsdp_iam_org.hospital1.id
}
```

```hcl
output "hsdp_iam_email_templates_ids" {
   value = data.hsdp_iam_email_templates.all.ids
}
```

## Argument Reference

* `organization_id` - (Required) The organization ID of the templates
* `locale` - (Optional) The locale to filter on

## Attributes Reference

The following attributes are exported:

* `ids` - The IDs of the templates
* `types` - The types of the templates
* `formats` - The formats of the templates
* `locales` - The locales of the templates
* `messages` - The message bodies of the templates
* `from` - The From: header addresses of the templates
* `subjects` - The subjects of the templates
* `links` - The links of the templates
