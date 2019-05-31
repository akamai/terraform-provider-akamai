---
layout: "akamai"
page_title: "Akamai: property rules"
sidebar_current: "docs-akamai-resource-property-rules"
description: |-
  Property Rules
---

# akamai_property_rules



The `akamai_property_rules` provides the resource for configuring a nested block of property rules, criteria, and behaviors. A property’s main functionality is encapsulated in its set of rules and rules are composed of the matches and the behavior that applies under those matches.


## Example Usage

Basic usage:

```hcl
resource "akamai_property_rules" "terraform-demo-rules" {

rules {   
  {
    “name”: “Origin”
    "children": []
    “behaviors”: [
    {
    “name”: “origin”
    option {
           key = "originType"
           value = "CUSTOMER"
           }
    option {
           key = "hostname"
           value = "api.example.org”
           }
    option {
           key = "forwardHostHeader"
           value = "REQUEST_HOST_HEADER”
           }
    option {
           key = "cachekeyHostname"
           value = "ORIGIN_HOSTNAME”
           }
    option {
           key = "compress"
           value = true
           }
    option {
           key = "enableTrueClientIp"
           value = false
           }
    option {
           key = "httpPort"
           value = 80
           }
           }
   ],
   "criteria": [
    {
      "name": "path",
      option {
         key = "matchCaseSensitive"
         value = false
      }
      option {
        key = "matchOperator"
        value = "IS_ONE_OF"
      }
      option {
      key = "values"
      value = [
      "/api/*"
      ]
      }

   }
     ]
"criteriaMustSatisfy": "all",
"comments": “Direct to origin based for APIs based on the path ."

}

{
       "name": "APIGateway"
   "children": []
  "behaviors": [
    {
      "name": "rapid"
    option {
        key = "enabled"
      value = true
           }
    }
    ]
}

}

```

## Argument Reference

The following arguments are supported:

The `rule` block supports:

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
