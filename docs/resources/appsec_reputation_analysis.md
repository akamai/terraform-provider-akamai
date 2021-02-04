---
layout: "akamai"
page_title: "Akamai: ReputationAnalysis"
subcategory: "Application Security"
description: |-
  ReputationAnalysis
---

# resource_akamai_appsec_reputation_analysis

The `resource_akamai_appsec_reputation_analysis` resource allows you to toggle the reputation analysis settings for a given security policy. The `forward_to_http_header` parameter indicates whether to add client reputation details to requests forwarded to origin in an HTTP header. The `forward_shared_ip_to_http_header_siem` parameter indicates whether to add value indicating that shared IPs are included in HTTP header and SIEM integration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to set reputation analysis settings
resource "akamai_appsec_reputation_analysis" "reputation_analysis" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  forward_to_http_header = true
  forward_shared_ip_to_http_header_siem = true
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security_policy_id to which the settings should be applied.

* `forward_to_http_header` - (Required) Whether to add client reputation details to requests forwarded to origin in an HTTP header.

* `forward_shared_ip_to_http_header_siem` - (Required) Whether to add value indicating that shared IPs are included in HTTP header and SIEM integration.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

