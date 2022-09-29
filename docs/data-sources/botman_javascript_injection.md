---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_javascript_injection (Beta)

**Scopes**: Security policy

Returns information about the JavaScript injection rules assigned to a security policy.

Use the [akamai_botman_javascript_injection](../resources/akamai_botman_javascript_injection) resource to modify your existing JavaScript injection rules.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/javascript-injection](https://techdocs.akamai.com/bot-manager/reference/get-javascript-injection-rules). Returns information about all the JavaScript injection rules associated with a security policy.

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

data "akamai_botman_javascript_injection" "javascript_injection" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "javascript_injection_json" {
  value = data.akamai_botman_javascript_injection.javascript_injection.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the JavaScript injection rules.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the JavaScript injection rules.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about the JavaScript injection rules assigned to a security policy.

**See also**:

- [Set up JavaScript injection](https://techdocs.akamai.com/bot-manager/docs/set-up-javascript-injection)
