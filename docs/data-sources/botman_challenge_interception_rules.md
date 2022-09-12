---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_challenge_interception_rules

**Scopes**: Security configuration

Returns information about the challenge interception rules available in the specified security configuration. Challenge interception rules help ensure that challenge actions work with AJAX (Asynchronous JavaScript and XML) requests. 

To modify an existing challenge interception rule, use the [akamai_botman_challenge_interception_rules](../resources/akamai_botman_challenge_interception_rules) resource.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/challenge-interception-rules](https://techdocs.akamai.com/bot-manager/reference/get-challenge-interception-rules). Returns information about your challenge interception rules.

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

data "akamai_botman_challenge_interception_rules" "interception_rules" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "interception_rules_json" {
  value = data.akamai_botman_challenge_interception_rules.interception_rules.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the challenge interception rules.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your challenge interception rules.
