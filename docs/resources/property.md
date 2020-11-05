---
layout: "akamai"
page_title: "Akamai: property"
subcategory: "Provisioning"
description: |-
  Create and update Akamai Properties
---

# akamai_property

The `akamai_property` resource represents an Akamai property configuration, allowing you to create,
update, and activate properties on the Akamai platform. NOTE: in 0.10 and earlier version this resource also 
controlled cpcode, origin, rules, and hostname associations but the logic was broken out into individual resources 
starting with 1.0.

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
    rules       = data.akamai_rules_template.example.json
}
```

## Argument Reference

The following arguments are supported:

### Property Basics

* `name` — (Required) The property name.
* `contact` — (Required) One or more email addresses to inform about activation changes.
* `contract_id` — (Required) The Contract ID.  Can be provided with or without `ctr_` prefix.
* `group_id` — (Required) The Group ID. Can be provided with or without `grp_` prefix.
* `product_id` — (Required) The Product ID. Can be provided with or without `prd_` prefix.
* `hostnames` — (Required) A map of public hostnames to edge hostnames (e.g. `{"example.org" = "example.org.edgesuite.net"}`)
* `rules` — (Required) A JSON encoded rule tree for given property. This should be provided in a form of complete json rule tree (see: [`akamai_property_rules`](../data-sources/property_rules.html))
* `rule_format` — (Optional) The rule format to use ([more](https://developer.akamai.com/api/core_features/property_manager/v1.html#getruleformats)) if not provided then the latest version will be used.

### Deprecated
* `contract` — (Deprecated) synonym of contract_id for legacy purposes
* `group` — (Deprecated) synonym of group_id for legacy purposes
* `product` — (Deprecated) synonym of product_id for legacy purposes

### Attribute Reference

The following attributes are returned:

* `warnings` — The contents of `warnings` field returned by the API.
* `errors` — The contents of `errors` field returned by the API.
* `latest_version` — The version of property on which the rules are created/updated - provider always uses latest or creates a new version if latest is not editable.
* `production_version` — the current version of the property active on the production network.
* `staging_version` — the current version of the property active on the staging network.
