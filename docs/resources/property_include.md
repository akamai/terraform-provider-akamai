---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_include (Beta)

Use the `akamai_property_include` resource to create an include and its rule tree.

## Basic usage

```hcl
resource "akamai_property_include" "my_example" {
    contract_id = "ctr_1-AB123"
    group_id    = "grp_12345"
    product_id  = "prd_123456"
    name        = "my.new.include.com"
    rule_format = "v2022-10-18"
    type        = "MICROSERVICES"
}
```

## Argument reference

This resource supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the optional `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the optional `grp_` prefix.
* `product_id` - (Required) A product's unique ID, including the `prd_` prefix. See [Common Product IDs](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/shared-resources#common-product-ids) for more information.
* `name` - (Required) The descriptive name for the include.
* `rules` - (Optional) Include's rules as JSON.
* `rule_format` - (Required) Indicates the versioned set of features and criteria. See [Rule format schemas](https://techdocs.akamai.com/property-mgr/reference/rule-format-schemas) to learn more.
* `type` - (Required) Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`. Use this field for filtering. `MICROSERVICES` allow different teams to work independently on different parts of a single site. `COMMON_SETTINGS` includes are useful for configurations that share a large number of settings, often managed by a central team.

## Attributes reference

This resource returns these attributes:

* `rule_errors` - Rule's validation errors. You need to resolve returned errors, as they block an activation.
* `rule_warnings` - Rule's validation warnings. You can activate a version that yields less severe warnings.
* `latest_version` - Returns the most recent version of the include.
* `staging_version` - The include version currently activated on the staging network.
* `production_version` - The include version currently activated on the production network.
