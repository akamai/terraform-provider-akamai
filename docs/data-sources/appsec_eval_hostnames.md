---
layout: "akamai"
page_title: "Akamai: EvalHostnames"
subcategory: "Application Security"
description: |-
 EvalHostnames
---

# akamai_appsec_eval_hostnames

The `akamai_appsec_eval_hostnames` data source allows you to retrieve the evaluation hostnames for a configuration version. Evaluation mode for hostnames is only available for Web Application Protector. Run hostnames in evaluation mode to see how your configuration settings protect traffic for that hostname before adding a hostname directly to a live configuration. An evaluation period lasts four weeks unless you stop the evaluation. Once you begin, the hostnames you evaluate start responding to traffic as if they are your current hostnames. However, instead of taking an action the evaluation hostnames log which action they would have taken if they were your actively-protected hostnames and not a test.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the hosts which are under evaluation in a config version

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
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

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `hostnames` - A list of the evaluation hostnames.

* `json` - A JSON-formatted list of the evaluation hostnames.

* `output_text` - A tabular display showing the evaluation hostnames.

