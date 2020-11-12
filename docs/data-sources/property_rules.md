---
layout: "akamai"
page_title: "Akamai: akamai_property_rules"
subcategory: "Provisioning"
description: |-
 Property ruletree
---

# akamai_property_rules


Use `akamai_property_rules` data source to query and retrieve the instance information and rule tree of an 
existing property instance.  allows searching across contracts and groups you may have access to.

## Basic Usage

Given a contract and group return what properties exist for the user:

datasource-example.tf
```hcl-terraform
datasource "akamai_property_rules" "my-example" {
    property_id = "prp_123"
    group_id = "grp_12345"
    contract_id = "ctr_1-AB123"
    version   = 3
}

output "property_match" {
  value = data.akamai_property_rules.my-example
}
```

## Argument Reference

The following arguments are supported:

* `contract_id` — (Required) The Contract ID.  Can be provided with or without `ctr_` prefix.
* `group_id` — (Required) The Group ID. Can be provided with or without `grp_` prefix.
* `property_id` — (Required) The property ID.  Can be provided with or without `prp_` prefix.
* `version` — (Optional) The version to return. (default: latest)

## Attributes Reference

The following are the return attributes:

* `json` — PAPIs response to the query.
