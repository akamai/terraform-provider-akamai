---
layout: akamai
subcategory: Application Security
---

## akamai_appsec_security_policy

**Scopes**: Security configuration; security policy

Returns information about your security policies.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies](https://techdocs.akamai.com/application-security/reference/get-policies)

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

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_security_policy" "security_policies" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "security_policies_list" {
  value = data.akamai_appsec_security_policy.security_policies.security_policy_id_list
}

output "security_policies_json" {
  value = data.akamai_appsec_security_policy.security_policies.json
}

output "security_policies_text" {
  value = data.akamai_appsec_security_policy.security_policies.output_text
}

data "akamai_appsec_security_policy" "specific_security_policy" {
  config_id            = data.akamai_appsec_configuration.configuration.config_id
  security_policy_name = "APIs"
}

output "specific_security_policy_id" {
  value = data.akamai_appsec_security_policy.specific_security_policy.security_policy_id
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the security policies.
- `security_policy_name`. (Optional). Name of the security policy you want to return information for (be sure to reference the policy name and not the policy ID). If not included, information is returned for all your security policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of the security policy information.
- `output_text`. Tabular report showing the ID and name of all your security policies.
- `security_policy_id`. ID of the security policy. Included only if the `security_policy_name` argument is included in your Terraform configuration file.
- `security_policy_id_list`. List of all your security policy IDs.
