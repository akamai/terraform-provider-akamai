---
layout: "akamai"
page_title: "Akamai: property activation"
subcategory: "Provisioning"
description: |-
  Property Activation
---

# akamai_property_activation

The `akamai_property_activation` provides the resource for activating a property in the appropriate environment. Once 
you are satisfied with any version of a property, an activation deploys it, either to the Akamai staging or production 
network. You activate a specific version, but the same version can be activated separately more than once.  
It is a best practice to have production activation depend on staging so if any problems are detected in staging 
the change does not progress to production.

## Example Usage

Basic usage:

```hcl
locals {
     email = "user@example.org"
     rule_format = "v2020-03-04"
}

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
     # manually specifying version allows production to lag behind staging until qualified by testing on staging urls.
     version = 3 
     # manually declairing a dependency on staging means production activation will not update if staging update fails -  
     # useful when both target same version.  The example does not depict this approach. However, this practice is 
     # recommended even when you edit production version by hand as shown in this example.
     depends_on = [
        akamai_property_activation.example_staging
     ]
     contact  = [local.email] 
}
```

## Argument Reference

The following arguments are supported:

* `property_id` — (Required) The property ID.  Can be provided with or without `prp_` prefix.
* `contact` — (Required) One or more email addresses to inform about activation changes.
* `network` — (Optional) Akamai network to activate on. Allowed values `STAGING` or `PRODUCTION` (Default: `STAGING`).
* `version` — (Required) The version to activate. Note: this field used to be optional but now depends on property to identify latest instead of calculating it locally.  This association helps keep the dependency tree properly aligned. 

### Deprecated Arguments
* `property` — (Deprecated) synonym of property_id for legacy purposes

## Attribute Reference

The following attributes are returned:

* `id` - unique identifier for this activation
* `warnings` - any warnings which may arise when the operation is executed by the infrastructure.
* `errors` - any errors which may arise when the operation is executed by the infrastructure
* `activation_id` - activation ID while an activation is in progress.
* `status` - the current activation status