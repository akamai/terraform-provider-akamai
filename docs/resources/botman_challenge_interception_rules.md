---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_challenge_interception_rules

**Scopes**: Security configuration

Modifies a challenge interception rule. 

To configure challenge interception rules you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `challenge_interception_rules` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To view your current challenge interception rules use the [akamai_botman_challenge_interception_rules](../data-sources/akamai_botman_challenge_interception_rules) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/response-actions/challenge-interception-rules](https://techdocs.akamai.com/bot-manager/reference/put-challenge-interception-rules). Updates the specified challenge interception rule.

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

resource "akamai_botman_challenge_interception_rules" "interception_rules" {
  config_id                    = data.akamai_appsec_configuration.configuration.config_id
  challenge_interception_rules = file("${path.module}/challenge_interception_rules.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the challenge interception rules
- `challenge_interception_rules` (Required). JSON-formatted collection of challenge interception rule settings and values. In the preceding sample code, the syntax `file("${path.module}/challenge_interception_rules.json")` points to the location of a JSON file containing the challenge interception rules settings and values.
