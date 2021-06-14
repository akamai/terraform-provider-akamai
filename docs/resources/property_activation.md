---
layout: "akamai"
page_title: "Akamai: property activation"
subcategory: "Provisioning"
description: |-
  Property Activation
---

# akamai_property_activation

~> **Note** Version 1.0.0 of the Akamai Terraform Provider is now available for the Property Provisioning module. To upgrade to this version, you have to update this resource. See the [Upgrade to Version 1.0.0](../guides/1.0_migration.md) for details.

The `akamai_property_activation` resource lets you activate a property version. An activation deploys the version to either the Akamai staging or production network. You can activate a specific version multiple times if you need to.  

Before activating on production, activate on staging first. This way you can detect any problems in staging before your changes progress to production.


## Example usage

Basic usage:

```hcl
locals {
     email = "user@example.org"
     rule_format = "v2020-03-04"
}

resource "akamai_property" "example" {
    name    = "terraform-demo"
    product_id  = "prd_SPM"
    contract_id = var.contractid
    group_id    = var.groupid
    hostnames = {
       "example.org" = "example.org.edgesuite.net"
       "www.example.org" = "example.org.edgesuite.net"
       "sub.example.org" = "sub.example.org.edgesuite.net"
    }
    rule_format = local.rule_format
    # line below here is assumed to be defined but left out for example brevity
    rules       = file("${path.module}/main.json")
}

resource "akamai_property_activation" "example_staging" {
     property_id = akamai_property.example.id
     contact  = [local.email]
     # NOTE: Specifying a version as shown here will target the latest version created. This latest version will always be activated in staging.
     version  = akamai_property.example.latest_version
     # not specifying network will target STAGING
}

resource "akamai_property_activation" "example_prod" {
     property_id = akamai_property.example.id
     network  = "PRODUCTION"
     # manually specifying version allows production to lag behind staging until qualified by testing on staging URLs.
     version = 3
     # manually declaring a dependency on staging means production activation will not update if staging update fails -  
     # useful when both target same version.  The example does not depict this approach. However, this practice is
     # recommended even when you edit production version by hand as shown in this example.
     depends_on = [
        akamai_property_activation.example_staging
     ]
     contact  = [local.email]
}
```

## Argument reference

The following arguments are supported:

* `property_id` - (Required) The property’s unique identifier, including the `prp_` prefix.
* `contact` - (Required) One or more email addresses to send activation status changes to.
* `version` - (Required) The property version to activate. Previously this field was optional. It now depends on the `akamai_property` resource to identify latest instead of calculating it locally.  This association helps keep the dependency tree properly aligned. To always use the latest version, enter this value `{resource}.{resource identifier}.{field name}`. Using the example code above, the entry would be `akamai_property.example.latest_version` since we want the value of the `latest_version` attribute in the `akamai_property` resource labeled `example`.
* `network` - (Optional) Akamai network to activate on, either `STAGING` or `PRODUCTION`. `STAGING` is the default.
* `auto_acknowledge_rule_warnings` - (Optional) Whether the activation should proceed despite any warnings. By default set to `true`.

### Deprecated arguments

* `property` - (Deprecated) Replaced by `property_id`. Maintained for legacy purposes.

## Attribute reference

The following attributes are returned:

* `id` - The unique identifier for this activation.
* `warnings` - The contents of `warnings` field returned by the API. For more information see [Errors](https://developer.akamai.com/api/core_features/property_manager/v1.html#errors) in the PAPI documentation.
* `errors` - The contents of `errors` field returned by the API. For more information see [Errors](https://developer.akamai.com/api/core_features/property_manager/v1.html#errors) in the PAPI documentation.
* `activation_id` - The ID given to the activation event while it's in progress.
* `status` - The property version’s activation status on the selected network.

### Deprecated attributes

* `rule_warnings` - (Deprecated) Rule warnings are no longer maintained in the state file. You can still see the warnings in logs.
