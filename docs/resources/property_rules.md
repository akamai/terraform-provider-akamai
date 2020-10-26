---
layout: "akamai"
page_title: "Akamai: property rules"
subcategory: "Provisioning"
description: |-
  Create and update Akamai property rule tree
---

# akamai_property

The `akamai_property_rules` resource represents an Akamai property rule tree configuration, allowing you to create and
update rule tree for given property.

## Example Usage

Basic usage:

```hcl
resource "akamai_property_rules" "example" {
    contract_id = var.contractid
    group_id = var.groupid
    property_id = var.propertyid
    rules       = <<-EOF
        {
           "name": "default",
                  "behaviors": [
            ...
        }
EOF
}
```

## Argument Reference

The following arguments are supported:

### Argument reference

* `contract_id` — (Required) The contract ID. Can be provided with or without `ctr_` prefix.
* `group_id` — (Required) The group ID. Can be provided with or without `grp_` prefix.
* `property_id` — (Required) The property ID. Can be provided with or without `prp_` prefix.
* `rules` — (Required) The rule tree for given property. This should be provided in a form of complete json rule tree.
* `rule_format` — (Optional) This parameter is currently not utilized - `latest` will always be used.

### Attribute Reference

The following attributes are returned:

* `version` — The version of property on which the rules are created/updated - provider always uses latest or creates a new version if latest is not editable.
* `warnings` — The contents of `warnings` field returned by the API.
* `errors` — The contents of `errors` field returned by the API.