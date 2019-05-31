---
layout: "akamai"
page_title: "Akamai: property"
sidebar_current: "docs-akamai-resource-property"
description: |-
  Create and update Akamai Properties
---

# akamai_property

The `akamai_property` resource represents an Akamai property configuration, allowing you to create,
update, and activate properties on the Akamai platform. 

## Example Usage

Basic usage:

```hcl
resource "akamai_property" "terraform-demo-web" {
    name    = "terraform-demo"
    contact = ["user@example.org"]

    product  = "prd_SPM"
    contract = "contract_####"
    group    = "grp_####"
    cp_code  = "cpc_#####"

    edge_hostname_map = "${merge(akamai_secure_edge_hostname.terraform-demo.edgehostmap)}"

    rule_format = "v2018-02-27"
    rules       = "${data.local_file.terraform-demo.content}"
    variables   = "${akamai_property_variables.origin.json}"
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
