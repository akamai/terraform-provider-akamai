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
resource "akamai_property_activation" "example" {
     name     = "${akamai_property.example.name}"
     contact  = ["user@example.org"] 
     hostname =  ["example.org"]
     contract = "ctr_####"
     group    = "grp_###"
     network  = "STAGING"
     activate = "${var.akamai_property_activate}"
}
```

## Argument Reference

The following arguments are supported:

* `property` — (Required) The property ID.
* `version` — (Optional) The version to activate (Default: latest)
* `network` — (Optional) Akamai network to activate on. Allowed values `staging` or `production` (Default: `staging`).
* `activate` — (Optional, boolean) Whether to activate the property on the network. (Default: `true`).
* `contact` — (Required) One or more email addresses to inform about activation changes.

## Attribute Reference

The follwing attributes are returned:

* `status` — the current activation status