---
layout: "akamai"
page_title: "Akamai: property"
subcategory: "Provisioning"
description: |-
  Create and update Akamai properties.
---

# akamai_property

~> **Note** Version 1.0.0 of the Akamai Terraform Provider is now available for the Property Provisioning module. To upgrade to this version, you have to update this resource. See the [Upgrade to Version 1.0.0](../guides/1.0_migration.md) for details.

The `akamai_property` resource represents an Akamai property configuration.
This resource lets you to create, update, and activate properties on the
Akamai platform.

Akamai’s edge network caches your web assets near to servers that request them.
A property provides the main way to control how edge servers respond to various
kinds of requests for those assets. Properties apply rules to a set of hostnames,
and you can only apply one property at a time to any given hostname. Each property
is assigned to a product, which determines which behaviors you can use. Each
property’s default rule needs a valid content provider (CP) code assigned to bill
and report for the service.

~> **Note** In version 0.10 and earlier of this resource, it also controlled content provider (CP) codes, origin settings, rules, and hostname associations. Starting with version 1.0.0, this logic is broken out into individual resources.

## Example usage

Basic usage:

```hcl
resource "akamai_property" "example" {
    name    = "terraform-demo"
    product_id  = "prd_SPM"
    contract_id = var.contractid
    group_id    = var.groupid
    hostnames {                                     # Hostname configuration
      cname_from = "example.com"
      cname_to = "example.com.edgekey.net"
      cert_provisioning_type = "DEFAULT"
    }
    hostnames {
      cname_from = "www.example.com"
      cname_to = "example.com.edgesuite.net"
      cert_provisioning_type = "CPS_MANAGED"
    }
    rule_format = "v2020-03-04"
    rules       = data.akamai_property_rules_template.example.json
}
```

## Argument reference

This resource supports these arguments:

* `name` - (Required) The property name.
* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the `grp_` prefix.
* `product_id` - (Required to create, otherwise Optional) A product's unique ID, including the `prd_` prefix.
* `hostnames` - (Optional) A mapping of public hostnames to edge hostnames. See the [`akamai_property_hostnames`](../data-sources/property_hostnames.md) data source for details on the necessary DNS configuration.

    ~> **Note** Starting from version 1.5.0, the `hostnames` argument supports a new block type. If you created your code and state in version 1.4 or earlier, you need to manually update your configuration and replace the previous input for `hostnames` with the new syntax. This error indicates that the state is outdated: `Error: missing expected [`. To fix it, remove `akamai_property` from the state and import it again.

    Requires these additional arguments:

      * `cname_from` - (Required) A string containing the original origin's hostname. For example, `"example.org"`.
      * `cname_to` - (Required) A string containing the hostname for edge content. For example,  `"example.org.edgesuite.net"`.
      * `cert_provisioning_type` - (Required) The certificate’s provisioning type, either the default `CPS_MANAGED` type for the custom certificates you provision with the [Certificate Provisioning System (CPS)](https://learn.akamai.com/en-us/products/core_features/certificate_provisioning_system.html), or `DEFAULT` for certificates provisioned automatically.
* `rules` - (Optional) A JSON-encoded rule tree for a given property. For this argument, you need to enter a complete JSON rule tree, unless you set up a series of JSON templates. See the [`akamai_property_rules`](../data-sources/property_rules.md) data source.
* `rule_format` - (Optional) The [rule format](https://developer.akamai.com/api/core_features/property_manager/v1.html#getruleformats) to use. Uses the latest rule format by default.

### Deprecated arguments

* `contract` - (Deprecated) Replaced by `contract_id`. Maintained for legacy purposes.
* `group` - (Deprecated) Replaced by `group_id`. Maintained for legacy purposes.
* `product` - (Deprecated) Optional argument replaced by the now required `product_id`. Maintained for legacy purposes.

## Attribute reference

The resource returns these attributes:

* `rule_errors` - The contents of `errors` field returned by the API. For more information see [Errors](https://developer.akamai.com/api/core_features/property_manager/v1.html#errors) in the PAPI documentation.
* `latest_version` - The version of the property you've created or updated rules for. The Akamai Provider always uses the latest version or creates a new version if latest is not editable.
* `production_version` - The current version of the property active on the Akamai production network.
* `staging_version` - The current version of the property active on the Akamai staging network.

### Deprecated attributes

* `rule_warnings` - (Deprecated) Rule warnings are no longer maintained in the state file. You can still see the warnings in logs.

## Import

Basic Usage:

```hcl
resource "akamai_property" "example" {
    # (resource arguments)
  }
```

You can import the latest Akamai property version by using either the `property_id` or a comma-delimited
string of the property, contract, and group IDs. You'll need to enter the string of IDs if the property belongs to multiple groups or contracts.

If using the string of IDs, you need to enter them in this order:

`property_id,contract_id,group_id`

To import a specific property version, pass additional parameters, either:

* `LATEST` to import the latest version of the property, regardless of whether it's active or not. This works the same as providing just the `property_id` or a string of the property, contract, and group IDs, which is the default behavior.
* `PRODUCTION`, `PROD`, or `P` to import the latest version activated on the production environment.
* `STAGING`, `STAGE`, `STAG`, or `S` to import the latest version activated on the staging environment.
* Version number or version number with the `ver_` prefix to import a specific property version. For example `3` and `ver_3` correspond to the same version number.


Here are some examples for the latest property version:

```shell
$ terraform import akamai_property.example prp_123
```

Or

```shell
$ terraform import akamai_property.example prp_123,ctr_1-AB123,grp_123
```

Here are some examples for the latest active property version on the production network:

```shell
$ terraform import akamai_property.example prp_123,P
```

Or

```shell
$ terraform import akamai_property.example prp_123,ctr_1-AB123,grp_123,PROD
```

Here are some examples for the specific property version:

```shell
$ terraform import akamai_property.example prp_123,3
```

Or

```shell
$ terraform import akamai_property.example prp_123,ctr_1-AB123,grp_123,ver_3
```
