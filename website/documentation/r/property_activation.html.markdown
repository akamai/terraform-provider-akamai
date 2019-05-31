---
layout: "akamai"
page_title: "Akamai: property activation"
sidebar_current: "docs-akamai-resource-property-activation"
description: |-
  Property Activation
---

# akamai_property_activation



The `akamai_property_activation` provides the resource for activating a property in the appropriate environment. Once you are satisfied with any version of a property, an activation deploys it, either to the Akamai staging or production network. You activate a specific version, but the same version can be activated separately more than once.


## Example Usage

Basic usage:

```hcl
resource "akamai_property_activation" "terraform-demo" {

     name     = "${akamai_property.a_sample_property.akamaideveloper.netname}"
     contact  = ["devrel@akamai.com"] 
     hostname =  ["akamaideveloper.net"]
     contract = "ctr_XXX"
     group    = "grp_XXX"
     network  = "STAGING"
     activate = "${var.akamai_property_activate}"

}



```

## Argument Reference

The following arguments are supported:

* `contract` — (Optional) The contract ID.
* `group` — (Optional) The group ID.
* `network` — (Optional) Akamai network to activate on. Allowed values staging (default) or production.
* `activate` — (Optional, boolean) Whether to activate the property on the network. Default: true.
* `name` — (Required) The property name.
* `hostname` — (Required) One or more public hostnames.
* `contact` — (Required) One or more email addresses to inform about activation changes.
* `account` — (Required) The account ID.
