---
layout: "akamai"
page_title: "Akamai: ReputationProfileAnalysis"
subcategory: "Application Security"
description: |-
 Reputation Profile Analysis
---

# akamai_appsec_reputation_profile_analysis

**Scopes**: Security policy

Returns information about the following two reputation analysis settings:

- `forwardToHTTPHeader`. When enabled, client reputation information associated with a request is forwarded to origin servers by using an HTTP header.
- `forwardSharedIPToHTTPHeaderAndSIEM`. When enabled, both the HTTP header and SIEM integration events include a value indicating that the IP addresses is shared address.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/reputation-analysis](https://techdocs.akamai.com/application-security/reference/get-reputation-analysis)

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

// USE CASE: User wants to view all the reputation analysis associated with a security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_reputation_profile_analysis" "reputation_analysis" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "reputation_analysis_text" {
  value = data.akamai_appsec_reputation_profile_analysis.reputation_analysis.output_text
}

output "reputation_analysis_json" {
  value = data.akamai_appsec_reputation_profile_analysis.reputation_analysis.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the reputation profile analysis settings.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the reputation profile analysis settings.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of the reputation analysis settings.
- `output_text`. Tabular report showing the reputation analysis settings.