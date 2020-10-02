---
layout: "akamai"
page_title: "Akamai: property rules"
subcategory: "Provisioning"
description: |-
  Property Rules
---

# akamai_property_rules

The `akamai_property_rules` data source allows you to configure a nested block of property rules, criteria, and behaviors. A property’s main functionality is encapsulated in its set of rules and rules are composed of the matches and the behavior that applies under those matches.

## Example Usage

Basic usage:

```hcl
data "akamai_property_rules" "example" {
  rules { # Default rule
  
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
}

resource "akamai_property" "example" {
  rules = "${data.akamai_property_rules.example.json}"
  
  // ...
}
```

## Argument Reference

The following arguments are supported:

The `rule` block supports:

* `is_secure` — (Optional) Whether the property is a secure (Enhanced TLS) property or not (top-level only).
* `criteria` — (Optional) One or more criteria to match requests on.
* `behavior` — (Optional) One or more behaviors to apply to requests that match.
* `rule` — (Optional) Child rules (may be nested five levels deep).

The `criteria` block supports:

* `name` — (Required) The name of the criteria.
* `option` — (Optional) One or more options for the criteria.


The `behavior` block supports:

* `name` — (Required) The name of the behavior.
* `option` — (Optional) One or more options for the behavior.

The `option` block supports:

* `key` — (Required) The option name.
* `value` — (Optional) A single value for the option.
* `values` — (Optional) An array of values for the option.

One of `value` or `values` is required.

## Attributes Reference

The following are the return attributes:

* `json` — The resulting JSON rule tree