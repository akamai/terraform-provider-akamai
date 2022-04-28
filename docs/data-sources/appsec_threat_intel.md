---
layout: "akamai"
page_title: "Akamai: Threat Intelligence"
subcategory: "Application Security"
description: |-
 Threat Intelligence
---

# akamai_appsec_threat_intel

**Scopes**: Security policy

Returns threat intelligence settings for a security policy Note that this data source is only available to organizations running the Adaptive Security Engine (ASE) beta. For more information on ASE, please contact your Akamai representative.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules/threat-intel](https://techdocs.akamai.com/application-security/reference/get-rules-threat-intel)

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
data "akamai_appsec_threat_intel" "threat_intel" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}
output "threat_intel" {
  value = data.akamai_appsec_threat_intel.threat_intel.threat_intel
}

output "json" {
  value = data.akamai_appsec_threat_intel.threat_intel.json
}
output "output_text" {
  value = data.akamai_appsec_threat_intel.threat_intel.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the threat intelligence settings.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the threat intelligence settings.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `threat_intel`. Reports the threat Intelligence setting, either **on** or **off**.
- `json`. JSON-formatted threat intelligence report.
- `output_text`. Tabular report of the threat intelligence information.