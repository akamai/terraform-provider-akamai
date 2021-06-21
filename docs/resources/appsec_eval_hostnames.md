---
layout: "akamai"
page_title: "Akamai: EvalHostnames"
subcategory: "Application Security"
description: |-
  EvalHostnames
---

# resource_akamai_appsec_eval_hostnames

The `resource_akamai_appsec_eval_hostnames` resource allows you to update the list of hostnames you want to evaluate for a configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to specify the hostnames to evaluate
resource "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = var.hostnames
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `hostnames` - (Required) A list of evaluation hostnames to be used for the specified configuration version.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

