---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_activation

Use the `akamai_property_activation` data source to retrieve activation information for a property version on staging
or production network.

## Example usage

Basic usage:

```hcl
locals {
  email       = "user@example.org"
  rule_format = "v2022-10-18"
}

resource "akamai_property" "example" {
  name        = "terraform-demo"
  product_id  = "prd_SPM"
  contract_id = var.contractid
  group_id    = var.groupid
  hostnames {
    cname_to               = "www.example.com.edgekey.net"
    cname_from             = "www.example.com"
    cert_provisioning_type = "DEFAULT"
  }
  rule_format = local.rule_format
  # line below here is assumed to be defined but left out for example brevity
  rules       = file("${path.module}/main.json")
}

resource "akamai_property_activation" "example_staging" {
      property_id = akamai_property.example.id
      contact     = [local.email]
      # NOTE: Specifying a version as shown here will target the latest version created. This latest version will always be activated in staging.
      version     = akamai_property.example.latest_version
      # not specifying network will target STAGING
      note        = "Sample activation"
}

data "akamai_property_activation" "example_staging" {
      property_id = akamai_property.example.id
      version     = akamai_property.example.latest_version
}
```

## Argument reference

The following arguments are supported:

* `property_id` - (Required) The property's unique identifier, including optional `prp_` prefix.
* `version` - (Required) The activated property version. The value depends on the `akamai_property` resource to identify the latest activated version instead of calculating it locally. To always use the latest version, set the variable to identify the resource you want to use: `akamai_property.{resource identifier}.latest_version`.
* `network` - (Optional) Akamai network to check the activation, either `STAGING` or `PRODUCTION`. If not specified, this defaults to `STAGING`.

## Attribute reference

The following attributes are returned:

* `id` - The unique identifier for this activation.
* `warnings` - The contents of `warnings` field returned by the API. For more information
  see [Errors](https://techdocs.akamai.com/property-mgr/reference/api-errors) in the PAPI documentation.
* `errors` - The contents of `errors` field returned by the API. For more information
  see [Errors](https://techdocs.akamai.com/property-mgr/reference/api-errors) in the PAPI documentation.
* `activation_id` - The activation's unique identifier, including optional `atv_` prefix.
* `status` - The property version's activation status on the selected network.
* `contact` - The email addresses to notify about the activation status changes.
* `note` - Log message assigned to the activation request.
