---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_include_rules (Beta)

Use the `akamai_property_include_rules` data source to query and get an include's rules. This data source lets you search across the contracts and groups you have access to.

## Basic usage

This example returns the include's rule tree based on the specified contract, group, and include IDs:

```hcl
data "akamai_property_include_rules" "my_example" {
    contract_id = "ctr_1-AB123"
    group_id    = "grp_12345"
    include_id  = "inc_123456"
    version     = 3
}

output "my_example" {
  value = data.akamai_property_include_rules.my_example
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the optional `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the optional `grp_` prefix.
* `include_id` - (Required) An include's unique ID with the optional `inc_` prefix.
* `version` - (Required) The include version you want to view the rules for.

## Attributes reference

This data source returns these attributes:

* `rules` - Include's rules as JSON.
* `name` - The descriptive name for the include.
* `rule_errors` - Rule's validation errors. You need to resolve returned errors, as they block an activation.
* `rule_warnings` - Rule's validation warnings. You can activate a version that yields non-blocking warnings.
* `rule_format` - Indicates the versioned set of features and criteria that are currently applied to a rule tree. See [Rule format schemas](https://techdocs.akamai.com/property-mgr/reference/rule-format-schemas) to learn more.
* `type` - Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`. Use this field for filtering. `MICROSERVICES` allow different teams to work independently on different parts of a single site. `COMMON_SETTINGS` includes are useful for configurations that share a large number of settings, often managed by a central team.
