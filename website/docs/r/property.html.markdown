---
layout: "akamai"
page_title: "Akamai: property"
sidebar_current: "docs-akamai-resource-property"
description: |-
  Create and update Akamai Properties
---

# akamai_property

The `akamai_property` resource represents an Akamai property configuration, allowing you to create,
update, and activate properties on the Akamai Platform. 

## Example Usage

Basic usage:

```hcl
resource "akamai_property" "example" {
  name = "example.com"

  contact = ["dshafik@akamai.com"]

  account_id = "act_XXXXXXXXX"
  contract_id = "ctr_C-XXXXXX"
  group_id = "grp_XXXXXXX"
  cp_code = "123456"

  hostname = ["example.org"]

  origin {
    is_secure = false
    hostname = "origin.example.org"
    forward_hostname = "ORIGIN_HOSTNAME"
  }
}
```

## Configuring Property Rules

To accommodate the numerous rule criteria and behaviors possible, the `akamai_property` resource uses `option` blocks inside
a `criteria` or `behavior` block.

```hcl
rules { # Default ule
  behavior { # Downstream Cache behavior
    name = "downstreamCache"
    option { # behavior option
      key = "behavior"
      value = "TUNNEL_ORIGIN"
    }
  }

  rule { # "Performance" child rule
    name = "Performance"
    rule { # "JPEG Images" child rule 
      name = "JPEG Images"
      behavior { # Adaptive Image Compression behavior
        name = "adaptiveImageCompression"
        # Options
        option {
          key = "tier1MobileCompressionMethod"
          value = "COMPRESS"
        }
        option {
          key = "tier1MobileCompressionValue"
          value = "80"
        }
        option {
          key = "tier2MobileCompressionMethod"
          value = "COMPRESS"
        }
      }
    }
  }
}
``` 

Each `option` block comprises of a `key` and a corresponding `value` (single value) or `values` (array of values).

> **Note:** You may nest `rule` blocks up to five levels deep. 

## Argument Reference

The following arguments are supported:

* `account_id` — (Required) The account ID
* `contract_id` — (Optional) The contract ID
* `group_id` — (Optional) The group ID
* `product_id` — (Optional) The product ID
* `network` — (Optional) Akamai network to activate on. Allowed values `staging` (default) or `production`.
* `activate` — (Optional, boolean) Whether to activate the property on the `network`. Default: `true`. 
* `cp_code` — (Required) The CP Code to use (or create)
* `name` — (Required) The property name
* `version` — 
* `rule_format` — (Optional) The rule format to use ([more](https://developer.akamai.com/api/luna/papi/overview.html#versioning))
* `ipv6` —  (Optional) Whether the property should use IPv6 to origin
* `hostname` — (Required) One or more public hostnames
* `contact` — (Required) One or more email addresses to inform about activation changes
* `edge_hostname` — (Optional) One or more edge hostnames (must be <= to the number of public hostnames))
* `clone_from` — (Optional) A property to clone
  * `property_id` — (Required) The ID of the property to clone
  * `version` — (Optional) The version of the property configuration to clone from (default: latest)
  * `etag` — (Optional) An etag for the property configuration that validates it has not changed (useful when cloning from latest)
  * `copy_hostnames` — (Optional, boolean) Whether to copy the hostnames configuration from the original property (if you copy hostnames and activate the property, it will **replace** the original configuration on the network)
* `origin` — (Optional) The property origin (an origin _must_ be specified to activate a property, but may be defined in your `rules` block)
  * `is_secure` — (Required) Whether the property configuration should be deployed to the the secure (TLS) Akamai network
  * `hostname` — (Required) The origin hostname
  * `port` — (Optional) The origin port to connect to (default: `80`)
  * `forward_hostname` — (Optional) The value for the `Hostname` header sent to origin. (default: `ORIGIN_HOSTNAME`)
  * `cache_key_hostname` — (Optional) The hostname uses for the cache key. (default: `ORIGIN_HOSTNAME`)
  * `compress` — (Optional, boolean) Whether origin supports gzip compression (default: `false`)
  * `enable_true_client_ip` — (Optional, boolean) Whether the `X-True-Client-IP` header should be sent to origin (default: `false`) 
* `rules` — (Optional) A nested block of property rules, criteria, and behaviors
  * `behavior` — (Optional) One or more behaviors to apply by default (use one `behavior` block for each behavior)
  * `rule` — (Optional) Child rules
  
  
The `rule` block supports:

* `criteria` — (Optional) One or more criteria to match requests on
* `behavior` — (Optional) One or more behaviors to apply to requests that match
* `rule` — (Optional) Child rules (may be nested five levels deep)

The `criteria` block supports:

* `name` — (Required) The name of the criteria
* `option` — (Optional) One or more options for the criteria
  
  
The `behavior` block supports:

* `name` — (Required) The name of the behavior
* `option` — (Optional) One or more options for the behavior
  
The `option` block supports:

* `key` — (Required) The option name
* `value` — (Optional) A single value for the option
* `values` — (Optional) An array of values for the option

One of `value` or `values` is required.

For more details on available Criteria and Behaviors, see the [Criteria](https://developer.akamai.com/api/luna/papi/criteria.html) and
[Behavior](https://developer.akamai.com/api/luna/papi/behaviors.html) documentation. 
