---
layout: "akamai"
page_title: "Akamai: WAPSelectedHostnames"
subcategory: "Application Security"
description: |-
  WAPSelectedHostnames
---

# akamai_appsec_wap_selected_hostnames

**Scopes**: Security policy

Modifies the list of hostnames to be protected or evaluated under a security configuration and security policy.
Either the evaluated hostnames or the protected hostnames may be omitted from or may be specified as an empty array (i.e., no hosts are to be protected or evaluated) in your Terraform configuration file.
However, at least one non-empty list must be included in the Terraform configuration file.
This resource is available only for Web Application Protector (WAP) accounts. Note that WAP selected hostnames is currently in beta. Please contact your Akamai representative for more information.

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

resource "akamai_appsec_wap_selected_hostnames" "appsecwap_selectedhostnames" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  protected_hosts    = ["documentation.akamai.com"]
  evaluated_hosts    = ["training.akamai.com"]
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the hostnames being protected or evaluated.
- `security_policy_id` (Required). Unique identifier of the security policy responsible for protecting or evaluating the specified hosts.
- `protected_hostnames` (Optional). JSON array of the hostnames to be protected. You must use either this argument or the `evaluated_hostnames` argument.
- `evaluated_hostnames` (Optional). JSON array of the hostnames to be evaluated. You must use either this argument or the `protected_hostnames` argument.