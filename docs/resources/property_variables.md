---
layout: "akamai"
page_title: "Akamai: property variables"
subcategory: "docs-akamai-resource-property-variables"
description: |-
  Property Variables
---

# akamai_property_variables


The `akamai_property_variables` allows you to implement dynamic functionality. You can perform conditional logic based on the variable’s value, and catch any unforeseen errors that execute on the edge at runtime.

Typical uses for variables include:

* Simplify configurations by reducing the number of rules and behaviors.
* Improve self serviceability by replacing or extending advanced metadata.
* Automate redirects, forward path rewrites, HTTP header and cookie manipulation.
* Move origin functionality to the edge.


## Example Usage

Basic usage:

```hcl
resource "akamai_property_variables" "origin" {
  variables {
     variable {
        name        = "PMUSER_ORIGIN"
        value       = "origin.example.org"
        description = "Origin Hostname"
        hidden      = true
        sensitive   = true
     }
  }
}
```

## Argument Reference

The following arguments are supported:

The `variables` block may contain many `variable` blocks which support the following arguments:

* `name` — (Required) The name of the variable.
* `value` — (Required) The default value to assign to the variable
* `description` — (Optional) A human-readable description
* `hidden` — (Required) Whether to hide the variable when debugging requests
* `sensitive` — (Required) Whether to obscure the value when debugging requests