---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_endpoint_coverage_report

**Scopes**: Universal (all endpoint coverage reports); operation

Returns your bot transactional endpoint coverage reports. These reports relay information about the Bot Manager protections applied to your API resources.

Use the `operation_id` argument to limit the returned data to information about a specified API operation. 

**Related API Endpoints**:

- [/appsec/v1/bot-endpoint-coverage-report](https://techdocs.akamai.com/bot-manager/reference/get-bot-endpoint-coverage-reports). Returns all your endpoint coverage reports.
- [appsec/v1/configs/{configId}/versions/{versionNumber}/bot-endpoint-coverage-report](https://techdocs.akamai.com/bot-manager/reference/get-bot-endpoint-coverage-reports-config-version). Returns coverage reports for the specified operation.

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

// USE CASE: User wants to return information for all coverage reports

data "akamai_botman_bot_endpoint_coverage_report" "coverage_reports" {
}

output "coverage_reports_json" {
  value = data.akamai_botman_bot_endpoint_coverage.coverage_reports.json
}

// USE CASE: User only wants to return coverage reports for the operation with the ID e0f89bb0-77d5-46f7-979d-e204e6fdc5a5

data "akamai_botman_bot_endpoint_coverage_report" "coverage_reports" {
  operation_id = "e0f89bb0-77d5-46f7-979d-e204e6fdc5a5"
}

output "coverage_reports_json" {
  value = data.akamai_botman_bot_endpoint_coverage.coverage_reports.json
}

// USE CASE: User only wants to return coverage reports for the specified security configuration

data "akamai_botman_bot_endpoint_coverage_report" "coverage_reports" {
  config_id  =  data.akamai_appsec_configuration.configuration.config_id
}

output "coverage_reports_json" {
  value = data.akamai_botman_bot_endpoint_coverage.coverage_reports.json
}

// USE CASE: User wants to return coverage reports for operation e0f89bb0-77d5-46f7-979d-e204e6fdc5a5 in the specified security configuration

data "akamai_botman_bot_endpoint_coverage_report" "coverage_reports" {
  config_id    = data.akamai_appsec_configuration.configuration.config_id
  operation_id = "e0f89bb0-77d5-46f7-979d-e204e6fdc5a5"
}

output "coverage_reports_json" {
  value = data.akamai_botman_bot_endpoint_coverage.coverage_reports.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Optional). Unique identifier of the security configuration associated with the coverage reports you want to return.
- `operation_id` (Optional). Unique identifier of the API operation to be returned. If omitted, information is returned for all your API operations.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your endpoint coverage reports.

**See also**:

- [See bot activity by API resource purpose](https://techdocs.akamai.com/bot-manager/docs/see-bot-activity-api-operation)
