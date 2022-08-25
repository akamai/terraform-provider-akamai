---
layout: "akamai"
page_title: "Akamai: ContractsGroups"
subcategory: "Application Security"
description: |-
 ContractsGroups
---


# akamai_appsec_contracts_groups

**Scopes**: Contract; group

Returns information about the contracts and groups associated with your account. Among other things, this information is required to create a new security configuration and to return a list of the hostnames available for use in a security policy. 

**Related API Endpoint**: [/appsec/v1/contracts-groups](https://techdocs.akamai.com/application-security/reference/get-contracts-groups)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the contracts and groups associated with their account.

data "akamai_appsec_contracts_groups" "contracts_groups" {
  contractid = "5-2WA382"
  groupid    = 12198
}

// USE CASE: User wants to display returned data in a table.

output "contracts_groups_list" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.output_text
}

output "contracts_groups_json" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.json
}

//USE CASE: User wants to return all available contracts and contract groups.

output "contract_groups_default_contractid" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.default_contractid
}

output "contract_groups_default_groupid" {
  value = data.akamai_appsec_contracts_groups.contracts_groups.default_groupid
}
```

## Argument Reference

This data source supports the following arguments:

- `contractid` (Optional). Unique identifier of an Akamai contract. If not included, information is returned for all the Akamai contracts associated with your account.
- `groupid` (Optional). Unique identifier of a contract group. If not included, information is returned for all the groups associated with your account.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of contract and group information.
- `output_text`. Tabular report of contract and group information.
- `default_contractid`. Default contract ID for the specified contract and group.
- `default_groupid`. Default group ID for the specified contract and group.