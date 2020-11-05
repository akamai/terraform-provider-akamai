---
layout: "akamai"
page_title: "Akamai: akamai_property_groups"
subcategory: "Provisioning"
description: |-
 Property groups
---

# akamai_property_groups


Use `akamai_property_groups` data source to list groups associated with an edgerc API token. 

## Basic Usage

Return what groups exist for the user:

datasource-example.tf
```hcl-terraform
datasource "akamai_property_groups" "my-example" {
}

output "property_match" {
  value = data.akamai_property_groups.my-example
}
```

## Argument Reference

No arguments are supported:

## Attributes Reference

The following are the return attributes:

* `json` — PAPIs response to the query.

Example PAPI response is as follows:
```json
{
    "accountId": "act_1-9ZYX87",
    "accountName": "Example.com",
    "groups": {
        "items": [
            {
                "groupName": "Example.com-1-1ABC234",
                "groupId": "grp_12345",
                "contractIds": [
                    "ctr_1-1ABC234"
                ]
            },
            {
                "groupName": "Test",
                "groupId": "grp_23455",
                "parentGroupId": "grp_12345",
                "contractIds": [
                    "ctr_1-1ABC234"
                ]
            }
        ]
    }
}
