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
network. You activate a specific version, but the same version can be activated separately more than once. Also note 
that it is a best practice to have production activation depend on staging so if any problems are detected in staging 
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
}

resource "akamai_property_version" "example" {
     contract_id = var.contractid
     group_id    = var.groupid
     property_id = akamai_property.example.id
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
     # we suggest depending on property_version and not property to ensure changes in version are completed before activations occur
     property_id = akamai_property_version.example.property_id
     contact  = [local.email] 
     # NOTE: specifying a version this way will target latest verson created and will mean latest will always be activated in staging.
     version  = akamai_property_version.exmple.version
     # not specifying network will target STAGING
}

resource "akamai_property_activation" "example_prod" {
     # we suggest depending on property_version and not property to ensure changes in version are completed before 
     # activations occur
     property_id = akamai_property_version.example.property_id
     network  = "PRODUCTION"
     # manually specifying version allows production to lag behind staging until qualified by testing on staging urls.
     version = 3 
     # manually declairing a dependency on staging means prod ativation will not update if staging update fails - this 
     # is useful when both target same version.  (which is not done here but is a safe pattern to stick with even when 
     # you edit prod by hand whihc is what was done below)
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
* `version` — (Required) The version to activate. Note: this field used to be optional but now depends on property_version to identify latest instead of calculating it locally.  This association helps keep the dependency tree properly aligned. 

### Deprecated Arguments
* `property` — (Deprecated) synonym of property_id for legacy purposes

## Attribute Reference

The following attributes are returned:

* `id` - unique identifier for this activation
* `warnings` - any warnings which may arise from server side CRUD operations
* `errors` - any errors which may arise from server side CRUD operations
* `activation_id` - activation ID while an activation is in progress.
* `status` - the current activation status