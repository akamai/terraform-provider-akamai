---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_client_side_security (Beta)

**Scopes**: Security configuration

Returns the client-side security settings for a security configuration. These settings configure Bot Manager to work in different situations.

To modify your existing client-side security settings, use the [akamai_botman_client_side_security](../resources/akamai_botman_client_side_security) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/client-side-security](https://techdocs.akamai.com/bot-manager/reference/get-client-side-security). Returns your client-side security settings.

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

data "akamai_botman_client_side_security" "client_side_security" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "client_side_security_json" {
  value = data.akamai_botman_client_side_security.client_side_security.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the client-side security settings.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your client-side security settings and setting values.

**See also**:

- [Handle bot management cookies](https://techdocs.akamai.com/bot-manager/docs/handle-bot-mgmt-cookies)
