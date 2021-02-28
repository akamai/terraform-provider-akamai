---
layout: "akamai"
page_title: "Akamai: akamai_property_hostnames"
subcategory: "Provisioning"
description: |-
 Property hostnames
---

# akamai_property_hostnames

~> **Note** Version 1.0.0 of the Akamai Terraform Provider is now available for the Provisioning module. To upgrade to the new version, you have to update this data source. See the [migration guide](../guides/1.0_migration.md) for details. 

Use the `akamai_property_hostnames` data source to query and retrieve the hostnames of 
an existing property. It also displays statuses of certificates associated to each hostname. This data source lets you search across the contracts 
and groups you have access to.

## Basic usage

This example returns the hostnames of a property based on the selected contract and group:

```hcl
datasource "akamai_property_hostnames" "my-example" {
    property_id = "prp_123"
    group_id = "grp_12345"
    contract_id = "ctr_1-AB123"
}

output "property_hostnames" {
  value = data.akamai_property_hostnames.my-example.hostnames
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 
* `group_id` - (Required) A group's unique ID, including the `grp_` prefix.
* `property_id` - (Required) A property's unique ID, including the `prp_` prefix.

## Attributes reference

This data source returns these attributes:

* `hostnames` - A list of hostnames for the property, including:
  * `cname_type` - A string containing cname type value of hostname.
  * `edge_hostname_id` - A string containing the edge hostname id including the `ehn_ `prefix of hostname.
  * `cname_from` - A string containing cname_from value of hostname.
  * `cname_to` - A string containing cname_to value of hostname.
  * `cert_provisioning_type` - A string containing cert_provisioning_type value of hostname.
  * `cert_status` - A list of cert statuses if they exist, including:
    * `target` - A string containing target value in the certificate.
    * `hostname` - A string containing hostname value in the certificate.
    * `production_status` - A string containing production status of certificate.
    * `staging_status` - A string containing staging status of certificate.
  
 
