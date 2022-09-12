---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_client

**Scopes**: Security configuration; custom client

Creates or modifies a custom client. CTo configure a custom client you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `custom_client` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-clients](https://techdocs.akamai.com/bot-manager/reference/post-custom-client). Creates a custom client.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-clients/{customClientId}](https://techdocs.akamai.com/bot-manager/reference/put-custom-client). Updates an existing custom client.

To review your existing custom clients, use the [akamai_botman_custom_client](../data-sources/akamai_botman_custom_client) data source.

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

// USE CASE: User wants to create a custom client

resource "akamai_botman_custom_client" "custom_client" {
  config_id     = data.akamai_appsec_configuration.configuration.config_id
  custom_client = file("${path.module}/custom_client.json")
}

// USE CASE: User wants to modify the custom client with the ID 1a3fd673-b9ed-4d11-8c9a-26157419ec77

resource "akamai_botman_custom_client" "custom_client" {
  config_id     = data.akamai_appsec_configuration.configuration.config_id
  custom_client = file("${path.module}/custom_client.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom client.
- `custom_client` (Required). JSON collection of settings and setting values for the custom client.  In the preceding sample code, the syntax `file("${path.module}/custom_client.json")` points to the location of a JSON file containing the custom_client settings and values.
    See [Create a custom client](https://techdocs.akamai.com/bot-manager/reference/post-custom-client) for detailed information on the settings available to you when creating or modifying a custom client.

**See also**:

- [Define custom clients like mobile apps](https://techdocs.akamai.com/bot-manager/docs/define-custom-clients-like-mobile-apps)
