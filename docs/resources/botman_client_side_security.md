---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_client_side_security (Beta)

**Scopes**: Security configuration

Modifies client-side security settings for the specified security configuration. To configure client-side security settings you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `client_side_security` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your client-side security settings, use the [akamai_botman_client_side_security](../data-sources/akamai_botman_client_side_security) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/client-side-security](https://techdocs.akamai.com/bot-manager/reference/put-client-side-security). Updates the client-side security settings for a security configuration.

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

resource "akamai_botman_client_side_security" "client_side_security" {
  config_id            = data.akamai_appsec_configuration.configuration.config_id
  client_side_security = file("${path.module}/client_side_security.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the client-side security settings.
- `client_side_security` (Required). JSON-formatted collection of client-side security settings and their values. In the preceding sample code, the syntax `file("${path.module}/client_side_security.json")` points to the location of a JSON file containing the client-side security settings and values.

**See also**:

- [Handle bot management cookies](https://techdocs.akamai.com/bot-manager/docs/handle-bot-mgmt-cookies)
