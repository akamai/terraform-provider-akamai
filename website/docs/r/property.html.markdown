---
layout: "akamai"
page_title: "Akamai: property"
sidebar_current: "docs-akamai-resource-property-config"
description: |-
  Create and update Akamai Properties
---

# akamai_property

The `akamai_property` resource represents an Akamai property configuration, allowing you to create,
update, and activate properties on the Akamai platform. 

## Example Usage

Basic usage:

```hcl
resource "akamai_property" "example" {
    name    = "terraform-demo"
    contact = ["user@example.org"]

    product  = "prd_SPM"
    contract = "ctr_####"
    group    = "grp_####"
    cp_code  = "cpc_#####"

    hostnames = {
      "example.org" = "example.org.edgesuite.net"
      "www.example.org" = "example.org.edgesuite.net"
      "sub.example.org" = "sub.example.org.edgesuite.net"
    }

    rule_format = "v2018-02-27"
    rules       = "${data.local_file.terraform-demo.content}"
    variables   = "${akamai_property_variables.origin.json}"
}
```

## Argument Reference

The following arguments are supported:

### Property Basics

* `account` — (Required) The account ID.
* `contract` — (Optional) The contract ID.
* `group` — (Optional) The group ID.
* `product` — (Optional) The product ID. (Default: `prd_SPM` for Ion)
* `name` — (Required) The property name.
* `contact` — (Required) One or more email addresses to inform about activation changes.
* `hostnames` — (Required) A map of public hostnames to edge hostnames (e.g. `{"example.org" = "example.org.edgesuite.net"}`)
* `is_secure` — (Optional) Whether the property is a secure (Enhanced TLS) property or not.

### Property Rules

* `rules` — (Required) A JSON encoded string of property rules (see: [`akamai_property_rules`](/docs/providers/akamai/d/property_rules.html))
* `rule_format` — (Optional) The rule format to use ([more](https://developer.akamai.com/api/core_features/property_manager/v1.html#getruleformats)).

In addition the specifying the rule tree in it's entirety, you can also set the default CP Code and Origin explicitly. *This will override your JSON configuration*.

* `cp_code` — (Optional) The CP Code id or name to use (or create). Required unless a [cpCode behavior](https://developer.akamai.com/api/core_features/property_manager/vlatest.html#cpcode) is present in the default rule.
* `origin` — (Optional) The property origin (an origin must be specified to activate a property, but may be defined in your rules block).
  * `hostname` — (Required) The origin hostname.
  * `port` — (Optional) The origin port to connect to (default: 80).
  * `forward_hostname` — (Optional) The value for the Hostname header sent to origin. (default: `ORIGIN_HOSTNAME`).
  * `cache_key_hostname` — (Optional) The hostname uses for the cache key. (default: `ORIGIN_HOSTNAME`).
  * `compress` — (Optional, boolean) Whether origin supports gzip compression (default: `false`).
  * `enable_true_client_ip` — (Optional, boolean) Whether the X-True-Client-IP header should be sent to origin (default: `false`).

You can also define property manager variables. *This will override your JSON configuration*.

* `variables` — (Optional) A JSON encoded string of property manager variable definitions (see: [`akamai_property_variables`](/docs/providers/akamai/r/property_variables.html))

### Attribute Reference

The following attributes are returned:

* `account` — the Account ID under which the property is created.
* `version` — the current version of the property config.
* `production_version` — the current version of the property active on the production network.
* `staging_version` — the current version of the property active on the staging network.
* `edge_hostnames` — the final public hostname to edge hostname map