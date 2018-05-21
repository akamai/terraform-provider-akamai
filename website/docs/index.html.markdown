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
  edgerc = "/path/to/.edgerc"
  papi_section = "papi"
  fastdns_section = "dns"
}

# Create a Property
resource "akamai_property" "example_property" {
  name = "www.example.org"

  # ...
}
```

## Authentication

The Akamai provider uses the standard Akamai Edgegrid authentication configuration,
providing a path to an `.edgerc` INI file, with one or more credential sections for each
service. You can read more about the .edgerc file [here](https://developer.akamai.com/introduction/Conf_Client.html#edgercformat). 

You can also specify credential values using environment variables. Environment variables take precedence over the `.edgerc` file. 

### Using an .edgerc file

To use an `.edgerc` file, you should configure the provider to specify a path. By
default it will look in the current users home directory. 

Usage:

```hcl
provider "akamai" {
  edgerc = "/path/to/.edgerc"
}
```

It is possible to use separate credentials for each service, by specifying their section. By
default the `default` section is used.

```hcl
provider "akamai" {
  edgerc = "/path/to/.edgerc"
  papi_section = "papi"
  fastdns_section = "dns"
}
```

### Environment variables

To use environment variables, you must create four different variables. The name
of the variables follow the following convention, which includes an optional `SECTION` name:

- `AKAMAI[_SECTION]_HOST`
- `AKAMAI[_SECTION]_CLIENT_TOKEN`
- `AKAMAI[_SECTION]_CLIENT_SECRET`
- `AKAMAI[_SECTION]_ACCESS_TOKEN`

If you are using the `default` section, you may either specify `DEFAULT` or omit the section name.

Usage:

```
$ export AKAMAI_HOST=akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net
$ export AKAMAI_ACCESS_TOKEN=akab-access-token-xxx-xxxxxxxxxxxxxxxx
$ export AKAMAI_CLIENT_TOKEN=akab-client-token-xxx-xxxxxxxxxxxxxxxx
$ export AKAMAI_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=
$ terraform plan
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `edgerc` - (Optional) The location of the `.edgerc` file containing credentials. Default: `$HOME/.edgerc`
* `papi_section` — (Optional) The credential section to use for the Property Manager API (PAPI). Default: `default`.
* `fastdns_section` — (Optional) The credential section to use for the Config DNS API. Default: `default`.

