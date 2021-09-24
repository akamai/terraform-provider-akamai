---
layout: "akamai"
page_title: "Akamai: EvalHostnames"
subcategory: "Application Security"
description: |-
  EvalHostnames
---

# akamai_appsec_eval_hostnames

**Scopes**: Security configuration

Modifies the list of hostnames evaluated while a security configuration is in evaluation mode.
During evaluation mode, hosts take no action of any kind when responding to traffic.
Instead, these hosts simply maintain a record of the actions they *would* have taken if they had been responding to live traffic in your production network.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/selected-hostnames/eval-hostnames](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevaluationhostnames)

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

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

// USE CASE: User wants to specify the hostnames to be evaluated in evaluation mode.

resource "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = ["documentation.akamai.com", "training.akamai.com", "videos.akamai.com"]
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration in evaluation mode.
- `hostnames` (Required). JSON array of hostnames to be used in the evaluation process. Note that this list replaces your existing list of evaluation hosts.

