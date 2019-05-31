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
    edgerc = "~/.edgerc"
    papi_section = "papi"
    # .. one of the two below
    dns_section = "dns"
    dnsv2_section = "dns"
    cps_section = "cps"
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
edgerc       =  "~/.edgerc"
}

```

You should specify separate credentials for each service. You may use the same .edgerc section for multiple services. By default the default section is used.

```hcl
provider "akamai" {
    edgerc = "~/.edgerc" 
    papi_section = "papi"
    dns_section = "dns"
    cps_section = "cps"
}

```


```

## Argument Reference

The following arguments are supported in the `provider` block:

* `edgerc` - (Optional) The location of the `.edgerc` file containing credentials. Default: `$HOME/.edgerc`
* `papi_section` — (Optional) The credential section to use for the Property Manager API (PAPI). Default: `default`.
* `dns_section` — (Optional) The credential section to use for the Config DNS API. Default: `default`.
* `cps_section` — (Optional) The credential section to use for the Config CPS. Default: `default`.

