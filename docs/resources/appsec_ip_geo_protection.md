---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_ip_geo_protection

**Scopes**: Security policy

Enables or disables IP/Geo protection for the specified configuration and security policy. When enabled, this allows your firewall to allow (or to block) clients based on their IP address or their geographic location.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/ip-geo-firewall](https://techdocs.akamai.com/application-security/reference/put-policy-protections)

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

// USE CASE: User wants to enable or disable IP/Geo protection.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_ip_geo_protection" "protection" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  enabled            = true
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the IP/Geo protection settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the IP/Geo protection settings being modified.
- `enabled` (Required). Set to **true** to enable IP/Geo protection; set to **false** to disable IP/Geo protection.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the current protection settings.