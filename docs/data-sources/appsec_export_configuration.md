---
layout: "akamai"
page_title: "Akamai: ExportConfiguration"
subcategory: "Application Security"
description: |-
 ExportConfiguration
---


# akamai_appsec_export_configuration

**Scopes**: Security configuration and version

Returns comprehensive details about a security configuration, including rate policies, security policies, rules, hostnames, and match targets.

**Related API Endpoint**: [/appsec/v1/export/configs/{configId}/versions/{versionNumber}](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurationversionexport)

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

data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version   = data.akamai_appsec_configuration.configuration.latest_version
  search    = ["securityPolicies", "selectedHosts"]
}

output "json" {
  value = data.akamai_appsec_export_configuration.export.json
}

output "text" {
  value = data.akamai_appsec_export_configuration.export.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration you want to return information for.
- `version` (Required). Version number of the security configuration.
- `search` (Optional). JSON array of strings specifying the types of information to be retrieved. Allowed values include:
> - **AdvancedSettingsLogging**
> - **AdvancedSettingsPrefetch**
> - **ApiRequestConstraints**
> - **AttackGroup**
> - **AttackGroupConditionException**
> - **Eval**
> - **EvalRuleConditionException**
> - **CustomDeny**
> - **CustomRule**
> - **CustomRuleAction**
> - **IPGeoFirewall**
> - **MatchTarget**
> - **PenaltyBox**
> - **RatePolicy**
> - **RatePolicyAction**
> - **ReputationProfile**
> - **ReputationProfileAction**
> - **Rule**
> - **RuleConditionException**
> - **SecurityPolicy**
> - **SiemSettings**
> - **SlowPost**


## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. Complete set of information about the specified security configuration version in JSON format. Includes the types available for the `search` parameter as well as additional fields such as `createDate` and `createdBy`.
- `output_text`. Tabular report showing the types of data specified in the `search` parameter. Valid only if the `search` parameter references at least one type.

