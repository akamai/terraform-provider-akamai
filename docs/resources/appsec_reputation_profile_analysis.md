---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_reputation_profile_analysis

**Scopes**: Security policy

Modifies the reputation analysis settings for a security policy. These settings include the following:

- The `forward_to_http_header` parameter, which indicates whether client reputation details are added to requests forwarded to origin in an HTTP header.
- The `forward_shared_ip_to_http_header_siem` parameter, which specifies whether a value is added indicating that shared IPs addresses are included in HTTP headers and in SIEM integration events.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/reputation-analysis](https://techdocs.akamai.com/application-security/reference/put-reputation-analysis)

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

data "akamai_appsec_reputation_profile_analysis" "reputation_analysis" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id     = "gms1_134637"
  forward_to_http_header = true
}

output "reputation_analysis_text" {
  value = data.akamai_appsec_reputation_profile_analysis.reputation_analysis.output_text
}

output "reputation_analysis_json" {
  value = data.akamai_appsec_reputation_profile_analysis.reputation_analysis.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the reputation profile analysis settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the reputation profile analysis settings being modified.
- `forward_to_http_header` (Optional). Set to **true** to add client reputation details to requests forwarded to the origin server in an HTTP header; set to `false` to leave reputation details out of these requests.
- `forward_shared_ip_to_http_header_siem` (Optional). Set to **true** to add a value indicating that shared IPs are included in HTTP header and SIEM integration; set to **false** to omit this value.