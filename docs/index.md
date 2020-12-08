---
layout: "akamai"
page_title: "Provider: Akamai"
description: |-
  Learn about the Akamai Terraform Provider
---

# Akamai Terraform Provider

Use the Akamai Terraform Provider to manage and provision your Akamai
configurations in Terraform. You can use the Akamai Provider today for
your Property Manager, Application Security, Edge DNS, and Global
Traffic Management configurations.

!> Version 1.0.0 of the Akamai Terraform Provider is a major release that's currently available for the Provisioning module. Before upgrading, you need to make changes to some of your Provisioning resources and data sources. See the [migration guide](guides/1.0_migration.md) for details.

Last updated: December 2020.

## Migrate to the newest version

If you're using the Provisioning module, the latest major version of the Akamai Provider is now available. See [1.0.0 Migration Guide](guides/1.0_migration.md) for more information.

## Workflows

Here are the most common workflows for the Akamai Provider:

* **Set up the Provider the first time.** To do this, finish reviewing this guide, then go to [Get Started with the Akamai Terraform Provider](guides/get_started_provider.md). When setting up the Provider, you need to choose an [authentication method](guides/akamai_provider_auth.md), and decide whether to import existing Akamai configurations, or create new ones.
* **Add a new module to your existing Akamai Provider configuration.** If the Akamai  Provider is already set up, and you're adding a new module, read the guide for the module you're adding:

  * [Application Security](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec)
  * [DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
  * [Global Traffic Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
  * [Provisioning/Property Manager](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
*  **Update settings for an existing module.** Use the reference information for resource and data sources listed under each module, like DNS or Provisioning. You can find this documentation on the panel to the left.

## Manage changes to your Akamai configurations

When you're using the Akamai Provider, you need to keep your
Terraform configurations up to date with changes made using Akamai
APIs, CLIs, and Control Center. You should review your network management
processes and update them to include the Akamai Provider.

For example, before updating your Akamai Provider configurations, you may want to
 run `terraform plan` first. You'll likely receive warnings
and suggested changes. Once you fix any issues, you can run `terraform plan`
again and make sure everything is in sync.

### Migrate to the newest version of the Akamai Provider

<!-- This section is a placeholder for the migration guide being developed. Likely need some overview text and a link to a separate migration guide.-->

## Links to resources

Here are some links to resources that can help get you started with the
Akamai Terraform Provider.

### New to Akamai?

If you're new to Akamai, here are some links to help you get started:

* [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started)
* [Akamai Community site](https://community.akamai.com/customers/s/)

### New to Terraform?

If you're new to Terraform, here are some links you might find helpful:

* [A Terraform Tutorial: Download to Installation, to Using Terraform](https://www.terraform.io/downloads.html)
* [A Brief Primer on Terraform's Configuration Language](https://www.terraform.io/docs/configuration/index.html)
* [If you want to learn about Terraform modules](https://www.terraform.io/docs/modules/index.html)
* [Terraform Glossary](https://www.terraform.io/docs/glossary.html)

## Available guides

* [Frequently Asked Questions](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/faq)
* [Get Started with DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
* [Get Started with GTM Domain Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
* [Get Started with the Provisioning Module](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
* [Get Started with Property Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property) <!--Name may change-->
* [Appendix](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/appendix)
<!--Want to rename to something like "Common codes and formats". Any suggestions?-->
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
  config_section = "default"
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
* `config_section` - (Optional) The credential section to use for all edgegrid calls. Default: `default`.
* `property_section` - (Deprecated) The credential section to use for the Property Manager API (PAPI). Default: `default`.
* `dns_section` - (Deprecated) The credential section to use for the Config DNS API. Default: `default`.
* `gtm_section` - (Deprecated) The credential section to use for the Config GTM API. Default: `default`.

## Additional Authentication Method - Inline Credentials

Outside of referenecing a local .edgerc file, you can specify credentials inline for each service used. Currently the provider supports `property` (PAPI), `dns` and `gtm` services.

To specify credentials inline, use the `property` or `dns` block to define credentials.

#### Argument Reference

* `config` - (Optional) Provide credentials for Terraform provider
  * `host` - (Required) The credential hostname
  * `access_token` - (Required) The credential access_token
  * `client_token` - (Required) The credential client_token
  * `client_secret` - (Required) The credential client_secret
  * `max_body` - (Optional) The credential max body to sign (in bytes, Default: `131072`)
  * `account_key` - (Optional) Account switch key to manage multiple accounts
* `property` - (Optional, Deprecated) Synonym to provide credentials for the Terraform provider using legacy `property` tag name
  * `host` - (Required) The credential hostname
  * `access_token` - (Required) The credential access_token
  * `client_token` - (Required) The credential client_token
  * `client_secret` - (Required) The credential client_secret
  * `max_body` - (Optional) The credential max body to sign (in bytes, Default: `131072`)
  * `account_key` - (Optional) Account switch key to manage multiple accounts
* `dns` - (Optional, Deprecated) Synonym to provide credentials for the Terraform provider using legacy `dns` tag name
  * `host` - (Required) The credential hostname
  * `access_token` - (Required) The credential access_token
  * `client_token` - (Required) The credential client_token
  * `client_secret` - (Required) The credential client_secret
  * `max_body` - (Optional) The credential max body to sign (in bytes, Default: `131072`)
  * `account_key` - (Optional) Account switch key to manage multiple accounts
* `gtm` - (Optional) (Optional, Deprecated) Synonym to provide credentials for the Terraform provider using legacy `gtm` tag name
  * `host` - (Required) The credential hostname
  * `access_token` - (Required) The credential access_token
  * `client_token` - (Required) The credential client_token
  * `client_secret` - (Required) The credential client_secret
  * `max_body` - (Optional) The credential max body to sign (in bytes, Default: `131072`)
  * `account_key` - (Optional) Account switch key to manage multiple accounts

## Environment Variables

You can specify credential values using environment variables. Environment variables take precedence over the contents of the `.edgerc` file.

Create environment variables in the format:

`AKAMAI{_SECTION_NAME}_*`

For example, if you specify `config_section = "papi"` you would set the following ENV variables:

* `AKAMAI_PAPI_HOST`
* `AKAMAI_PAPI_ACCESS_TOKEN`
* `AKAMAI_PAPI_CLIENT_TOKEN`
* `AKAMAI_PAPI_CLIENT_SECRET`
* `AKAMAI_PAPI_MAX_BODY` (optional)
* `AKAMAI_PAPI_ACCOUNT_KEY` (optional)

If the section name is `default`, you can omit it, instead using:

* `AKAMAI_HOST`
* `AKAMAI_ACCESS_TOKEN`
* `AKAMAI_CLIENT_TOKEN`
* `AKAMAI_CLIENT_SECRET`
* `AKAMAI_MAX_BODY` (optional)
* `AKAMAI_ACCOUNT_KEY` (optional)

## Guides

* [Frequently Asked Questions](guides/faq.md)
* [Get Started with DNS Zone Administration](guides/get_started_dns_zone.md)
* [Get Started with GTM Domain Administration](guides/get_started_gtm_domain.md)
* [Get Started with Identity and Access Management](guides/get_started_iam.md)
* [Get Started with Property Management](guides/get_started_property.md)
* [Appendix](guides/appendix.md)
