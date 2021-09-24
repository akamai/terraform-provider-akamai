---
layout: "akamai"
page_title: "Akamai: EvalProtectHost"
subcategory: "Application Security"
description: |-
  EvalProtectHost
---

# akamai_appsec_eval_protect_host

**Scopes**: Security configuration

**Important**: This data source is deprecated and may be removed in a future release. You may use the `akamai_appsec_wap_selected_hostnames` resource instead.

Moves hostnames being evaluated to active protection. When you move a hostname from the evaluation hostnames list that host is added to your security policy as a protected hostname and is removed from the collection of hosts being evaluated.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/protect-eval-hostnames](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putmoveevaluationhostnamestoprotection)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: User wants to move the evaluation hosts to the protected hosts list.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
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

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration in evaluation mode.
- `hostnames` (Required). JSON array of the hostnames to be moved from the evaluation hostname list to the protected hostname list.

