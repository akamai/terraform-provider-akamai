---
layout: "akamai"
page_title: "Akamai: EvalHostnames"
subcategory: "Application Security"
description: |-
 EvalHostnames
---


# akamai_appsec_eval_hostnames

**Scopes**: Security configuration

**Important**: This data source is deprecated and may be removed in a future release. You may use the `akamai_appsec_wap_selected_hostnames` data source instead.

Returns the evaluation hostnames for a configuration. In evaluation mode, you use evaluation hosts to monitor how well your configuration settings protect host traffic. (Note that the evaluation host isn't actually protected, and the host takes no action other than recording the actions it would have taken had it been on the production network.)

Evaluation mode for hostnames is available only for organizations running Web Application Protector.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/selected-hostnames/eval-hostnames](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevaluationhostnames)

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

// USE CASE: User wants to view the hosts being evaluated in evaluation mode.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "eval_hostnames" {
  value = data.akamai_appsec_eval_hostnames.eval_hostnames.hostnames
}

output "eval_hostnames_output" {
  value = data.akamai_appsec_eval_hostnames.eval_hostnames.output_text
}

output "eval_hostnames_json" {
  value = data.akamai_appsec_eval_hostnames.eval_hostnames.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration running in evaluation mode.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `hostnames`. List of evaluation hostnames.
- `json`. JSON-formatted list of evaluation hostnames.
- `output_text`. Tabular report showing evaluation hostnames.

