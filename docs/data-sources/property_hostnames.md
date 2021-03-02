---
layout: "akamai"
page_title: "Akamai: akamai_property_hostnames"
subcategory: "Provisioning"
description: |-
 Property hostnames
---

# akamai_property_hostnames

~> **Note** Version 1.0.0 of the Akamai Terraform Provider is now available for the Provisioning module. To upgrade to the new version, you have to update this data source. See the [migration guide](../guides/1.0_migration.md) for details.

Use the `akamai_property_hostnames` data source to query and retrieve hostnames and their certificate statuses for an existing property. This data source lets you search across the contracts and groups you have access to.

## Basic usage

This example returns the property's hostnames based on the selected contract and group:

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
  * `cname_type` - A string containing the hostname's cname type value.
  * `edge_hostname_id` - The edge hostname's unique ID, including the `ehn_` prefix.
  * `cname_from` - A string containing the original origin's hostname.
  * `cname_to` - A string containing the hostname for edge content.
  * `cert_provisioning_type` - The certificate’s provisioning type, either the default `CPS_MANAGED` type for the custom certificates you provision with the Certificate Provisioning System (CPS), or `DEFAULT` for certificates provisioned automatically.
  * `cert_status` - If applicable, this shows a list of certificate statuses, including:
    * `target` - The destination part of the CNAME record used to validate the certificate’s domain.
    * `hostname` - The hostname part of the CNAME record used to validate the certificate’s domain.
    * `production_status` - A string containing the status of the certificate deployment on the production network.
    * `staging_status` - A string containing the status of the certificate deployment on the staging network.

## Domain validation for DEFAULT certificates

If your `cert_provisioning_type = "DEFAULT"`, you need to perform domain validation to prove to the certificate authority that you control the domain and are authorized to create certificates for it.

In your DNS configuration, create a CNAME record and map the `cert_status.hostname` value to the `cert_status.target` value.
