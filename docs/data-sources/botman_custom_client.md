---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_client (Beta)

**Scopes**: Security configuration; custom client

Returns information about your custom clients.

Use the `custom_client_id` argument to limit returned data to the specified custom client.

To create or modify a custom client, use the [akamai_botman_custom_client](../resources/akamai_botman_custom_client) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-clients](https://techdocs.akamai.com/bot-manager/reference/get-custom-clients). Returns information about all your custom clients.
- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-clients/{customClientId}](https://techdocs.akamai.com/bot-manager/reference/get-custom-client). Returns information about the specified custom client.

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

// USE CASE: User wants to return information for all custom clients

data "akamai_botman_custom_client" "custom_clients" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}
output "custom_clients_json" {
  value = data.akamai_botman_custom_client.custom_clients.json
}

// USE CASE: User only wants to return information for the custom client with the ID 592b2f92-df19-4dd8-863a-7f305ca2c3c7

data "akamai_botman_custom_client" "custom_client" {
  config_id        = data.akamai_appsec_configuration.configuration.config_id
  custom_client_id = "592b2f92-df19-4dd8-863a-7f305ca2c3c7"
}

output "custom_client_json" {
  value = data.akamai_botman_custom_client.custom_client.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom clients.
- `custom_client_id` (Optional). Unique identifier of the custom client you want returned.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your custom clients. The returned information includes such things as the client type and the client platform (e.g., Android or iOS).

**See also**:

- [Define custom clients like mobile apps](https://techdocs.akamai.com/bot-manager/docs/define-custom-clients-like-mobile-apps)
