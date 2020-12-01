---
layout: "akamai"
page_title: "Akamai: Authenticate the Akamai Terraform Provider"
description: |-
  Learn how to set up authentication for the Akamai Terraform Provider.
---

# Authenticate the Akamai Terraform Provider
<!--Not sure about the name of this doc. -->

Authentication of Terraform configurations relies on the Akamai EdgeGrid
authentication scheme. The Akamai Terraform Provider code acts as a
wrapper for our APIs and reuses the same authentication mechanism. We
recommend storing your API credentials in a local .edgerc file.

See [Get Started with APIs](https://developer.akamai.com/api/getting-started) 
for more information about the process of creating Akamai credentials.

## Get authenticated 

The permissions you need for the Akamai Terraform Provider depends on
the subset of Akamai resources and data sources you'll use. Without
these permissions, your Terraform configurations won't execute.

To get authenticated you need to:

* Set up your API clients

* Add your local .edgerc file to your Akamai Terraform config

## Set up your API clients

Before you [create an API client](https://developer.akamai.com/api/getting-started#createanapiclient),
you need to:

1.  Determine which Akamai products you'll be using with Terraform.
2.  Find the API service name for the products you'll be adding.
    If you already have the API clients you need, you can add the credential
	to your local .edgerc file.
	<!--Need to go back to this one.-->

For example, if you're adding your Akamai properties to your existing
Terraform configuration, you'll need read-write permission to the
[Property Manager
API](https://developer.akamai.com/api/core_features/property_manager/v1.html).
In this case, you'll need to create an API client for the **Property
Manager (PAPI)** service, then add it to your .edgerc file.

Here's a list of the Akamai products available on Terraform and their
supporting APIs:

| **Product** | **API service name** |
|-------------|----------------------|
| Property Manager (Provisioning and Common modules) | Property Manager (PAPI)|
| Edge DNS (DNS) | DNS-Zone Record Management |
| Global Traffic Management | Traffic Management Configurations |
| Application Security | Application Security |

Once you create the supporting API clients you can update your local
`.edgerc file`.

### Add your local .edgerc file to your Akamai Terraform config


To reference a local .edgerc file, you add this line to the top of the
Akamai Terraform configuration file (akamai.tf): edgerc =
\"\~/.edgerc\".

The \~/.edgerc is the location of your file on your local machine. In
your Terraform files you can reference individual sections inside the
.edgerc file:

### Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    akamai = {
      source  = "hashicorp/akamai"
      version = "~> 0.10.0"
    }
  }
}

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


### Argument Reference

Arguments supported in the `provider` block:

* edgerc - (Optional) The location of the `.edgerc` file containing credentials. The default is `\$HOME/.edgerc`.
* config_section - (Optional) The credential section to use within the `.edgerc` file for all Edge Grid calls. If not included, uses credentials in the default section of the `.edgerc` file.

#### Deprecated Arguments

* property_section - (Deprecated) The credential section to use for the [Property Manager API](https://developer.akamai.com/api/core_features/property_manager/v1.html).
If not added, uses credentials in the default section of the `.edgerc` file.
* dns_section - (Deprecated) The credential section to use for the [Edge DNS Zone Management API](https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html). If not added, uses credentials in the default section of the `.edgerc` file.
* [gtm_section] - (Deprecated) The credential section to use for the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html). If not added, uses credentials in the default section of the `.edgerc` file.

## How Else Can I Authenticate on Akamai?

Referencing a local `.edgerc` file in your `akamai.tf` file is the preferred
authentication method. However there are a few other methods you can use
if necessary.

Authenticate using inline credentials
-------------------------------------

If needed, you can specify credentials inline for each resource or data
source used. To do this, you'll need to add a `config` block that
includes authentication information.

### Example Usage

<!--**NEED CODE EXAMPLE HERE.**-->

### Argument Reference

* `config` - (Optional) Provide credentials for Akamai provider. This block supports these arguments:
<!--Need better descriptions here. -->
  * `host` - (Required) The credential hostname.
  * `access_token` - (Required) The service's `access_token` from the `.edgerc` file.
  * `client_token` - (Required) The service's `client_token` from the `.edgerc` file.
  * `client_secret` - (Required) The service's `client_secret` from the `.edgerc` file.
  * `max_body` - (Optional) The service's `max_body` to sign in bytes. The default is 131072 bytes.
  * `account_key` - (Optional) If managing multiple accounts, the account ID you want to use when running Terraform commands. The account selected persists for all commands until you change it.

#### Deprecated Arguments

* `dns` - (Deprecated) Legacy Edge DNS API service argument for inline authentication. Used same arguments as the current `config` block.
* `gtm` - (Deprecated) Legacy Global Traffic Management API service argument for inline authentication. Used same arguments as the current `config` block.
* `property` - (Deprecated) Legacy Property Manager API service argument for inline authentication. Used same arguments as the current `config` block.

### Authenticate using environment variables

CHECK THAT THIS IS THE LATEST.

You can also use environment variables to set credential values.
Environment variables take precedence over the contents of the `.edgerc`
file.

Your environment variables should be in this format: `AKAMAI{_SECTION_NAME}_\*`

For example, if you're setting up the Provisioning module, you'll need to add a `config_section` block with these environment variables for your Property Manager API client:
<!--Go back and reread the above.-->

* `AKAMAI_PAPI_HOST`
* `AKAMAI_PAPI_ACCESS_TOKEN`
* `AKAMAI_PAPI_CLIENT_TOKEN`
* `AKAMAI_PAPI_CLIENT_SECRET`
* `AKAMAI_PAPI_MAX_BODY` (Optional)
* `AKAMAI_PAPI_ACCOUNT_KEY` (Optional)

If you're setting up variables for your `default` credentials, you can use these variables:

* `AKAMAI_HOST`
* `AKAMAI_ACCESS_TOKEN`
* `AKAMAI_CLIENT_TOKEN`
* `AKAMAI_CLIENT_SECRET`
* `AKAMAI_MAX_BODY` (Optional)
* `AKAMAI_ACCOUNT_KEY` (Optional)

### Example Usage

<!--**NEED CODE EXAMPLE HERE.**-->

### Argument Reference

<!--**Not sure if needed.**-->