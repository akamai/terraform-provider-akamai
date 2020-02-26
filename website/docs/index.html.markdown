---
layout: "akamai"
page_title: "Provider: Akamai"
sidebar_current: "docs-akamai-index"
description: |-
  Akamai
---

# Akamai Provider

The Akamai provider is used to interact with the Akamai platform for content
delivery, security, and performance.

To use this provider you must create API credentials valid for each service you want to
use. Learn more [here](https://developer.akamai.com/introduction/Prov_Creds.html).

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Akamai Provider
provider "akamai" {
    property {
      host = "${var.akamai_host}"
      access_token = "${var.akamai_access_token}"
      client_token = "${var.akamai_client_token}"
      client_secret = "${var.akamai_client_secret}"
    }
}


# Create a Property
resource "akamai_property" "example_property" {
  name = "www.example.org"

  # ...
}
```

## Authentication

You must specify credentials for each service used. Currently the provider supports `property` (PAPI), `dns` and gtm services.

You may use either a the Akamai standard `.edgerc` file, or you can specify the credentials inline.

### Inline Credentials

To specify credentials inline, use the `property` or `dns` block to define credentials.

#### Argument Reference

* `property` — (Optional) Provide credentials for the Property Manager API (papi)
  * `host` — (Required) The credential hostname
  * `access_token` — (Required) The credential access_token
  * `client_token` — (Required) The credential client_token
  * `client_secret` — (Required) The credential client_secret
  * `max_body` — (Optional) The credential max body to sign (in bytes, Default: `131072`)
* `dns` — (Optional) Provide credentials for the Edge DNS API (config-dns)
  * `host` — (Required) The credential hostname
  * `access_token` — (Required) The credential access_token
  * `client_token` — (Required) The credential client_token
  * `client_secret` — (Required) The credential client_secret
  * `max_body` — (Optional) The credential max body to sign (in bytes, Default: `131072`)
* `gtm` — (Optional) Provide credentials for the GTM Config API (config-gtm)
  * `host` — (Required) The credential hostname
  * `access_token` — (Required) The credential access_token
  * `client_token` — (Required) The credential client_token
  * `client_secret` — (Required) The credential client_secret
  * `max_body` — (Optional) The credential max body to sign (in bytes, Default: `131072`)

### Using an .edgerc file

The Akamai provider uses the standard Akamai Edgegrid authentication configuration,
providing a path to an `.edgerc` INI file, with one or more credential sections for each
service. You can read more about the .edgerc file [here](https://developer.akamai.com/introduction/Conf_Client.html#edgercformat).

To use an `.edgerc` file, you should configure the provider to specify a path. By
default it will look in the current users home directory.

Usage:

```hcl
provider "akamai" {
    edgerc =  "~/.edgerc"
}
```

You should specify separate credentials for each service. You may use the same .edgerc section for multiple services. By default the default section is used.

```hcl
provider "akamai" {
    edgerc = "~/.edgerc" 
    property_section = "papi"
    dns_section = "dns"
    gtm_section = "gtm"
}
```

#### Argument Reference

The following arguments are supported in the `provider` block:

* `edgerc` - (Optional) The location of the `.edgerc` file containing credentials. Default: `$HOME/.edgerc`
* `property_section` — (Optional) The credential section to use for the Property Manager API (PAPI). Default: `default`.
* `dns_section` — (Optional) The credential section to use for the Config DNS API. Default: `default`.
* `gtm_section` — (Optional) The credential section to use for the Config GTM API. Default: `default`.
* `cps_section` — (Optional) The credential section to use for the Config CPS. Default: `default`.

## Environment Variables

You can also specify credential values using environment variables. Environment variables take precedence over the contents of the `.edgerc` file.

Create environment variables in the format:

`AKAMAI{_SECTION_NAME}_*`

For example, if you specify `property_section = "papi"` you would set the following ENV variables:

* `AKAMAI_PAPI_HOST`
* `AKAMAI_PAPI_ACCESS_TOKEN`
* `AKAMAI_PAPI_CLIENT_TOKEN`
* `AKAMAI_PAPI_CLIENT_SECRET`
* `AKAMAI_PAPI_MAX_BODY` (optional)

If the section name is `default`, you can omit it, instead using:

* `AKAMAI_HOST`
* `AKAMAI_ACCESS_TOKEN`
* `AKAMAI_CLIENT_TOKEN`
* `AKAMAI_CLIENT_SECRET`
* `AKAMAI_MAX_BODY` (optional)
