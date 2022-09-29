---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_wap_selected_hostnames

**Scopes**: Security policy

Returns hostnames currently protected or being evaluated by a configuration and security policy.
This resource is available only to organizations running Web Application Protector (WAP).
Note that the WAP selected hostnames feature is currently in beta.
Please contact your Akamai representative for more information.

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

data "akamai_appsec_wap_selected_hostnames" "wap_selected_hostnames" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "protected_hostnames" {
  value = data.akamai_appsec_wap_selected_hostnames.wap_selected_hostnames.protected_hosts
}

output "evaluated_hostnames" {
  value = data.akamai_appsec_wap_selected_hostnames.wap_selected_hostnames.evaluated_hosts
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the hostnames.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the hostnames.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `protected_hostnames`. List of hostnames currently protected under the security configuration and security policy.
- `evaluated_hostnames`. List of hostnames currently being evaluated under the security configuration and security policy.
- `hostnames_json`. JSON-formatted report of the protected and evaluated hostnames.
- `output_text`. Tabular reports of the protected and evaluated hostnames.