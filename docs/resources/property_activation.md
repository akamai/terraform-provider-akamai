---
layout: "akamai"
page_title: "Akamai: property activation"
subcategory: "Provisioning"
description: |-
  Property Activation
---

# akamai_property_activation

The `akamai_property_activation` provides the resource for activating a property in the appropriate environment. Once you are satisfied with any version of a property, an activation deploys it, either to the Akamai staging or production network. You activate a specific version, but the same version can be activated separately more than once.

## Example Usage

Basic usage:

```hcl
resource "akamai_property_activation" "example" {
     property_id = "${akamai_property.example.id}"
     network  = "STAGING"
     activate = "${var.akamai_property_activate}"
     contact  = ["user@example.org"] 
}
```

## Argument Reference

The following arguments are supported:

* `property` - (Deprecated) The property.
* `property_id` - (Required) The property ID. Exclusive with `property`.
* `version` - (Optional) The version to activate. When unset it will activate the latest version of the property.
* `network` - (Optional) Akamai network to activate on. Allowed values `staging` or `production` (Default: `staging`).
* `activate` - (Deprecated, boolean) Whether to activate the property on the network. (Default: `true`).
* `contact` - (Required) One or more email addresses to inform about activation changes.

## Attribute Reference

The following attributes are returned:

* `id` - property ID + : + network (for example: `prp_1234:PRODUCTION`)
* `warnings` - any warnings which may arise from server side CRUD operations
* `errors` - any errors which may arise from server side CRUD operations
* `activation_id` - the activation ID
* `status` - the current activation status