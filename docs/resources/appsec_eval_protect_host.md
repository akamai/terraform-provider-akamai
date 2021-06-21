---
layout: "akamai"
page_title: "Akamai: EvalProtectHost"
subcategory: "Application Security"
description: |-
  EvalProtectHost
---

# resource_akamai_appsec_eval_protect_host

The `resource_akamai_appsec_eval_protect_host` resource allows you to move hostnames that you are evaluating to active protection. When you move a hostname from the evaluation hostnames list, it’s added to your security policy as a protected hostname.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to move the evaluation hosts to protected
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

resource "akamai_appsec_eval_protect_host" "protect_host" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = data.akamai_appsec_eval_hostnames.eval_hostnames.hostnames
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `hostnames` - (Required) The evaluation hostnames to be moved to active protection.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

