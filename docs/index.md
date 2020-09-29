---
layout: "akamai"
page_title: "Provider: Akamai"
subcategory: "docs-akamai-index"
description: |-
  Akamai
---

# Akamai Terraform Provider

The Akamai Terraform Provider is used to interact with Akamai and manage solutions for content
delivery, security, and performance.

Last updated: July 2020.

## Prerequisites

In order to get started with the Akamai Terraform Provider, the following is required:

* Create an Akamai API Client with the right permissions and valid credentials to authenticate your Akamai Terraform files. Your Akamai API Client will require read-write permissions to either the Property Manager API, DNS API and/or Traffic Management API depending on your use-case for using the Akamai Terraform Provider.

* Either import existing configurations with the [Akamai Terraform CLI](https://github.com/akamai/cli-terraform) or start from scratch with code examples. Note: Both Terraform and the Akamai Terraform CLI package come pre-installed in the Akamai Development Environment. Get more details in our [installation Instructions](https://developer.akamai.com/blog/2020/05/26/set-development-environment).

* Run the terraform init command to load the Akamai Terraform Provider.

* Run the terraform plan or terraform apply command to run your Terraform configuration.

Use the navigation to the left to read about the available resources available in the Akamai Terraform Provider.

## Authenticating Akamai Terraform configurations
Authentication of Terraform configurations relies on the Akamai EdgeGrid authentication scheme. The Akamai Terraform Provider code acts as a wrapper for our APIs and re-uses the same authentication mechanism. Note: We recommend storing your API credentials in a local .edgerc file.

Your Akamai API Client will require read-write permissions to either the Property Manager API, DNS API and/or Traffic Management API depending on your use-case for using the Akamai Terraform Provider. Note: Without these permissions, your Terraform configurations won’t execute.

The local .edgerc file can be referenced in the top of the Akamai Terraform configuration with edgerc = "~/.edgerc". Note: ~/.edgerc is the location of your file on your local machine. You are able to reference individual sections inside the .edgerc file by referencing papi_section = "default". Note: "default" is the name of the section stored in brackets in your .edgerc file.


## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    akamai = {
      source  = "hashicorp/akamai"
      version = "~> 0.9.1"
    }
  }
}

# Configure the Akamai Provider
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "papi"
  dns_section = "dns"
  gtm_section = "gtm"
}

# Create a Property
resource "akamai_property" "example_property" {
  name = "www.example.org"
  
  # ...
}

# Create a DNS Record
resource "akamai_dns_record" "example_record" {
  zone       = "example.org"
  name       = "www.example.org"
  recordtype = "CNAME"
  active     = true
  ttl        = 600
  target     = ["example.org.akamaized.net."]
}
```

Terraform 0.12 and earlier:

```hcl
# Configure the Akamai Provider
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "papi"
  dns_section = "dns"
  gtm_section = "gtm"
}

# Create a Property
resource "akamai_property" "example_property" {
  name = "www.example.org"
  
  # ...
}

# Create a DNS Record
resource "akamai_dns_record" "example_record" {
  zone       = "example.org"
  name       = "www.example.org"
  recordtype = "CNAME"
  active     = true
  ttl        = 600
  target     = ["example.org.akamaized.net."]
}

```

#### Argument Reference

The following arguments are supported in the `provider` block:

* `edgerc` - (Optional) The location of the `.edgerc` file containing credentials. Default: `$HOME/.edgerc`
* `property_section` — (Optional) The credential section to use for the Property Manager API (PAPI). Default: `default`.
* `dns_section` — (Optional) The credential section to use for the Config DNS API. Default: `default`.
* `gtm_section` — (Optional) The credential section to use for the Config GTM API. Default: `default`.

## Additional Authentication Method - Inline Credentials

Outside of referenecing a local .edgerc file, you can specify credentials inline for each service used. Currently the provider supports `property` (PAPI), `dns` and `gtm` services.

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

## Guides

* [Frequently Asked Questions](guides/faq.md)
* [Get Started with DNS Zone Administration](guides/get_started_dns_zone.md)
* [Get Started with GTM Domain Administration](guides/get_started_gtm_domain.md)
* [Get Started with Property Management](guides/get_started_property.md)
* [Appendix](guides/appendix.md)