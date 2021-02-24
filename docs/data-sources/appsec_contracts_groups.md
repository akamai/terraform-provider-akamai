---
layout: "akamai"
page_title: "Akamai: ContractsGroups"
subcategory: "Application Security"
description: |-
 ContractsGroups
---

# akamai_appsec_contracts_groups

Use the `akamai_appsec_contracts_groups` data source to retrieve information about the contracts and groups for your account. Each object contains the contract, groups associated with the contract, and whether Kona Site Defender or Web Application Protector is the product for that contract. You’ll need this information when you create a new security configuration or when you want to get a list of hostnames still available for use in a security policy. The information available via this data source is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcontractsandgroupswithksdorwaf).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to see contract group details in an account
data "akamai_appsec_contracts_groups" "contracts_groups" {
  contractid = var.contractid
  groupid = var.groupid
}

//tabular data of contractid, displayname and group id
output "contracts_groups_list" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.output_text
}

output "contracts_groups_json" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.json
}

//returns any of the contract/group
output "contract_groups_default_contractid" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.default_contractid
}

output "contract_groups_default_groupid" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.default_groupid
}
```

## Argument Reference

The following arguments are supported:

## Attributes Reference

* `contractid` - (Optional) The ID of a contract for which to retrieve information.

* `groupid` - (Optional) The ID of a group for which to retrieve information.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the contract and group information.

* `output_text` - A tabular display showing the contract and group information.

* `default_contractid` - The default contract ID for the specified contract and group.

* `default_groupid` - The default group ID for the specified contract and group.

