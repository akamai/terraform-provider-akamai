---
layout: "akamai"
page_title: "Akamai: Authenticate the Akamai Terraform Provider"
description: |-
  Learn how to set up authentication for the Akamai Terraform Provider.
---

# Authenticate the Akamai Terraform Provider

Authentication of Akamai Terraform Provider relies on the Akamai EdgeGrid
authentication scheme. The Akamai Provider code acts as a
wrapper for our APIs and reuses the same authentication mechanism. We
recommend storing your API credentials in a local `.edgerc` file.

See [Get Started with APIs](https://developer.akamai.com/api/getting-started) 
for more information about the process of creating Akamai credentials.

## Get authenticated 

The permissions you need for the Akamai Provider depends on
the subset of Akamai resources and data sources you'll use. Without
these permissions, your Terraform configurations won't execute.

To get authenticated you need to:

* Set up your API clients.
* Add your local .edgerc file to your Akamai Provider configuration.

## Set up your API clients

Before you [create an API client](https://developer.akamai.com/api/getting-started#createanapiclient),
you need to:

1.  Determine which Akamai Provider modules you'll be using with Terraform.
2.  Find the API service name for the modules you'll be adding.
3.  If you already have the API clients you need, you can add the credential
	to your local `.edgerc` file.

For example, if you're adding your Akamai properties to your existing
Terraform configuration, you'll be using the Provisioning module. For this module,
 you need to create an API client for the **Property Manager (PAPI)** service, 
 then add it to your `.edgerc` file.

Here's a list of the Akamai modules available on Terraform and their
supporting API service names:

| **Module** | **API service name** |
|-------------|----------------------|
| Property Manager (Provisioning and Common modules) | Property Manager (PAPI) |
| Edge DNS (DNS) | DNS Zone Management |
| Global Traffic Management | Global Traffic Management |
| Application Security | Application Security |

Once you create the supporting API clients you can update your local
`.edgerc` file.

## Add your local .edgerc file to your Akamai Provider config


To reference a local `.edgerc` file, you add this line to the top of the
Akamai Provider configuration file (`akamai.tf`): 

```
edgerc = "~/.edgerc"
```

`~/.edgerc` is the location of your file on your local machine. In
your Terraform files, you can reference individual sections inside the
`.edgerc` file:

### Example usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    akamai = {
      source  = "hashicorp/akamai"
    }
  }
}

# Configure the Akamai Provider
provider "akamai" {
  edgerc = "~/.edgerc"
  config_section = "default"
}

# Create a property
resource "akamai_property" "example_property" {
  name = "www.example.org"
  
  # ...
}

# Create a DNS record
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

# Create a property
resource "akamai_property" "example_property" {
  name = "www.example.org"
  
  # ...
}

# Create a DNS record
resource "akamai_dns_record" "example_record" {
  zone       = "example.org"
  name       = "www.example.org"
  recordtype = "CNAME"
  active     = true
  ttl        = 600
  target     = ["example.org.akamaized.net."]
}
```


### Argument reference

Arguments supported in the `provider` block:

* edgerc - (Optional) The location of the `.edgerc` file containing credentials. The default is `$HOME/.edgerc`.
* config_section - (Optional) The credential section to use within the `.edgerc` file for all EdgeGrid calls. If you don't use `config_section`, the Akamai Provider uses the credentials in the `default` section of the `.edgerc` file.

#### Deprecated arguments

* property_section - (Deprecated) The credential section to use for the [Property Manager API](https://developer.akamai.com/api/core_features/property_manager/v1.html).
If you don't use `property_section`, the Akamai Provider uses the credentials in the `default` section of the `.edgerc` file.
* dns_section - (Deprecated) The credential section to use for the [Edge DNS Zone Management API](https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html). If you don't use `dns_section`, the Akamai Provider uses the credentials in the `default` section of the `.edgerc` file.
* [gtm_section] - (Deprecated) The credential section to use for the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html). If you don't use `gtm_section`, the Akamai Provider uses the credentials in the `default` section of the `.edgerc` file.

## Authenticate using inline credentials

Referencing a local `.edgerc` file in your `akamai.tf` file is the preferred
authentication method. If needed, you can specify credentials inline for each 
resource or data source. To do this, add a `config` block that includes authentication information.

### Example usage
When using inline credentials, you need to update the Akamai Provider block to this:

```
provider "akamai" {
  config { # NOTE : Replace values with valid edge token values
    client_secret = "abcde7926304iHSUIYGE+9876XYZabcd/Io89="
    host = "akaa-blablablablablabla-blabla920937blabla.luna-dev.akamaiapis.net"
    access_token = "akaa-abcdef8398279-zyx90380397"
    client_token = "akaa-lksalhdsuw8993-982hj2kjdb2u"
  }
}
```

You don't need the `edgerc` or `config_section` attributes that you'd use if you were adding your local `.edgerc` file to your Akamai Provider configuration.

### Argument reference

* `config` - (Optional) Provide credentials for Akamai Provider. This block supports these arguments:
  * `host` - (Required) The base credential hostname without the protocol.
  * `access_token` - (Required) The service's access token from the `.edgerc` file.
  * `client_token` - (Required) The service's client token from the `.edgerc` file.
  * `client_secret` - (Required) The service's client secret from the `.edgerc` file.
  * `max_body` - (Optional) The service's maximum data payload size in bytes. 
  * `account_key` - (Optional) If managing multiple accounts, the account ID you want to use when running Terraform commands. The account selected persists for all commands until you change it.

#### Deprecated arguments

* `dns` - (Deprecated) Legacy Edge DNS API service argument for inline authentication. Used same arguments as the current `config` block.
* `gtm` - (Deprecated) Legacy Global Traffic Management API service argument for inline authentication. Used same arguments as the current `config` block.
* `property` - (Deprecated) Legacy Property Manager API service argument for inline authentication. Used same arguments as the current `config` block.

## Authenticate using environment variables

You can also use environment variables to set credential values.
Environment variables take precedence over the settings in the `.edgerc`
file.

Your environment variables should be in this format: `AKAMAI{_SECTION_NAME}_*`

For example, if you're setting up the Provisioning module, you'll need an API client for Property Manager. In your `akamai.tf` file, you'll need to add a `config_section` block with these environment variables:

* `AKAMAI_PAPI_HOST`
* `AKAMAI_PAPI_ACCESS_TOKEN`
* `AKAMAI_PAPI_CLIENT_TOKEN`
* `AKAMAI_PAPI_CLIENT_SECRET`
* `AKAMAI_PAPI_MAX_BODY` (Optional)
* `AKAMAI_PAPI_ACCOUNT_KEY` (Optional)

These variables cover the arguments you'd enter if using inline variables. 

If you're setting up variables for your `default` credentials, you can use these variables:

* `AKAMAI_HOST` 
* `AKAMAI_ACCESS_TOKEN` 
* `AKAMAI_CLIENT_TOKEN` 
* `AKAMAI_CLIENT_SECRET`
* `AKAMAI_MAX_BODY` (Optional) 
* `AKAMAI_ACCOUNT_KEY` (Optional) 

### Example usage

First, set up the base Akamai Provider block so that it's empty:

```
provider "akamai" {}
```

Then, in your terminal, set the values for your variables and run `terraform apply`: 

``` 
AKAMAI_HOST=akaa-hfwxy7qdv3a6v5pc-xupo52fjzb3yhpgw.luna-dev.akamaiapis.net \
AKAMAI_ACCESS_TOKEN=akaa-gzzvy3juqgjts7ao-xmhsyq4tsif5likt \
AKAMAI_CLIENT_TOKEN=akaa-ds7bmxvsl4rtjii6-vpuo5l5mf7n2z4bn \
AKAMAI_CLIENT_SECRET=vzy6SEo8nja23k2XXMYofnD10Xag+ju2iwwABu/QsOo= \
terraform apply
```

### Variable reference
When using variables, you'll need to set them up based on the sections of your `.edgerc` file they represent. Your environment variables should be in this format: `AKAMAI{_SECTION_NAME}_*`

These are the variables for the `default` section of your `.edgerc` and what they represent: 

* `AKAMAI_HOST` - (Required) The base credential hostname without the protocol.
* `AKAMAI_ACCESS_TOKEN` - (Required) The service's access token from the `.edgerc` file.
* `AKAMAI_CLIENT_TOKEN` - (Required) The service's client token from the `.edgerc` file.
* `AKAMAI_CLIENT_SECRET` - (Required) The service's client secret from the `.edgerc` file.
* `AKAMAI_MAX_BODY` - (Optional) The service's maximum data payload size in bytes.
* `AKAMAI_ACCOUNT_KEY` - (Optional) If managing multiple accounts, the account ID you want to use when running Terraform commands. The account selected persists for all commands until you change it.