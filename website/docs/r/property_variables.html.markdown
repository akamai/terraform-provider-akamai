---
layout: "akamai"
page_title: "Akamai: property variables"
sidebar_current: "docs-akamai-resource-property-variables"
description: |-
  Property Variables
---

# akamai_property_variables



The `akamai_property_variables` provides the resource that allows you to implement dynamic functionality .You can perform conditional logic based on the variable’s value, and catch any unforeseen errors that execute on the edge at runtime.


Typical uses for variables include:
<ol>
<li>Simplify configurations by reducing the number of rules and behaviors.</li>
<li>Improve self serviceability by replacing or extending advanced metadata.</li>
<li>Automate redirects, forward path rewrites, HTTP header and cookie manipulation.</li>
<li>Move origin functionality to the edge.</li>


## Example Usage

Basic usage:

```hcl
resource "akamai_property_variables" "origin" {
   variables {
       variable {
          name        = "PMUSER_ORIGIN"
          value       = “test-origin.akamaideveloper.net"
          description = "Terraform Demo Origin"
          hidden      = true
          sensitive   = false
                }
             }
}


```

## Argument Reference

The following arguments are supported:

* `account` — (Required) The account ID.
* `contract` — (Optional) The contract ID.
* `group` — (Optional) The group ID.
* `product` — (Optional) The product ID.
* `cp_code` — (Required) The CP Code id or name to use (or create).
* `name` — (Required) The property name.
* `rule_format` — (Optional) The rule format to use (more).
* `contact` — (Required) One or more email addresses to inform about activation changes.
* `edge_hostname` — (Optional) One or more edge hostnames (must be <= to the number of public hostnames)
* `edge_hostname_map` — (Optional) The edge hostname mapping.
* `origin` — (Optional) The property origin (an origin must be specified to activate a property, but may be defined in your rules block).
* `is_secure` — (Required) Whether the property configuration should be deployed to the the secure (TLS) Akamai network.
* `hostname` — (Required) The origin hostname.
* `port` — (Optional) The origin port to connect to (default: 80).
* `forward_hostname` — (Optional) The value for the Hostname header sent to origin. (default: ORIGIN_HOSTNAME).
* `cache_key_hostname` — (Optional) The hostname uses for the cache key. (default: ORIGIN_HOSTNAME).
* `compress` — (Optional, boolean) Whether origin supports gzip compression (default: false).
* `enable_true_client_ip` — (Optional, boolean) Whether the X-True-Client-IP header should be sent to origin (default: false).
* `rules` –  (Optional) The rules comprising of the matches and behavior as a json string.
* `name` — (Required) The variable name.
* `value` —(Required) The value of the variable name.
* `description` —(Optional) The description of the variable.
* `hidden` — (Required, boolean) Whether the variable is hidden.
* `sensitive` — (Required, boolean) Whether the variable is sensitive.
