---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_transactional_endpoint_protection (Beta)

**Scopes**: Security configuration

Updates transactional endpoint protection settings. To configure a transactional endpoint protection you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `transactional_endpoint_protection` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current protection settings, use the [akamai_botman_transactional_endpoint_protection](../data-sources/akamai_botman_transactional_endpoint_protection) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/transactional-endpoint-protection](https://techdocs.akamai.com/bot-manager/reference/put-transactional-endpoint-protection). Updates transactional endpoint protections.

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

resource "akamai_botman_transactional_endpoint_protection" "endpoint_protection" {
  config_id                         = data.akamai_appsec_configuration.configuration.config_id
  transactional_endpoint_protection = file("${path.module}/transactional_endpoint_protection.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the transactional endpoint protections being updated.
- `transactional_endpoint` (Required). JSON-formatted collection of transactional endpoint protection settings and setting values.  In the preceding sample code, the syntax `file("${path.module}/transactional_endpoint_protection.json")` points to the location of a JSON file containing the transactional endpoint protection settings and values.
