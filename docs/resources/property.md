---
layout: "akamai"
page_title: "Akamai: property"
subcategory: "Provisioning"
description: |-
  Create and update Akamai Properties
---

# akamai_property

The `akamai_property` resource represents an Akamai property configuration. This resources lets you to create,
update, and activate properties on the Akamai platform. 

Akamai’s edge network caches your web assets near to servers that request them. A property provides the main way to control how edge servers respond to various kinds of requests for those assets. Properties apply rules to a set of hostnames, and you can only apply one property at a time to any given hostname. Each property is assigned to a product, which determines which behaviors you can use. Each property’s default rule needs a valid content provider code (CP code) assigned to bill and report for the service.

> __NOTE:__ In version 0.10 and earlier of this resource, it also controlled 
content provider (CP) codes, origin settings, rules, and hostname associations. Starting with version 1.0, this logic is broken out into individual resources.

## Example Usage

Basic usage:

```hcl
resource "akamai_property" "example" {
    name    = "terraform-demo"
    contact = ["user@example.org"]
    product_id  = "prd_SPM"
    contract_id = var.contractid
    group_id    = var.groupid
    hostnames = {
      "example.org" = "example.org.edgesuite.net"
      "www.example.org" = "example.org.edgesuite.net" 
      "sub.example.org" = "sub.example.org.edgesuite.net"
    }
    rule_format = "v2020-03-04"
    rules       = data.akamai_property_rules_template.example.json
}
```

## Argument Reference

This resource supports these arguments:

* `name` - (Required) The property name.
* `contact` - (Required) One or more email addresses to send activation status changes to.
* `contract_id` - (Required) A contract's unique ID. If your ID doesn't include the `ctr_` prefix, the Akamai Provider appends it to your entry for processing purposes. 
* `group_id` - (Required) A group's unique ID. If your ID doesn't include the `grp_` prefix, the Akamai provider appends it to your entry for processing purposes.
* `product_id` - (Required) A product's unique ID. If your ID doesn't include the `prd_` prefix, the Akamai Provider appends it to your entry for processing purposes.
* `hostnames` - (Required) A mapping of public hostnames to edge hostnames (For example: `{"example.org" = "example.org.edgesuite.net"}`)
* `rules` - (Required) A JSON-encoded rule tree for a given property. For this argument, you need to enter a complete JSON rule tree, unless you set up a series of JSON templates. See the [`akamai_property_rules`](/docs/providers/akamai/d/property_rules.html))
* `rule_format` - (Optional) The [rule format](https://developer.akamai.com/api/core_features/property_manager/v1.html#getruleformats) to use. Uses the latest rule format by default.

### Deprecated Arguments

* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `group` - (Deprecated) Replaced by `group_id`. Maintained for legacy purposes.
* `product` - (Deprecated) Optional argument replaced by the now required `product_id`. Maintained for legacy purposes.

## Attribute Reference

The resource returns these attributes:

* `warnings` - The contents of `warnings` field returned by the API. For more information see [Errors](https://developer.akamai.com/api/core_features/property_manager/v1.html#errors) in the PAPI documentation.
* `errors` - The contents of `errors` field returned by the API. For more information see [Errors](https://developer.akamai.com/api/core_features/property_manager/v1.html#errors) in the PAPI documentation.
* `latest_version` - The version of the property you've created or updated rules for. The Akamai Provider always uses the latest version or creates a new version if latest is not editable.
* `production_version` - The current version of the property active on the Akamai production network.
* `staging_version` - The current version of the property active on the Akamai staging network.