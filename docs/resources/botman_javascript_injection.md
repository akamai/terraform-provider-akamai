---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_javascript_injection

**Scopes**: Security policy

Updates an existing JavaScript injection rule. To configure JavaScript injection rules you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `javascript_injection` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your existing JavaScript injection rules, use the [akamai_botman_javascript_injection](../data-sources/akamai_botman_javascript_injection) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/javascript-injection](https://techdocs.akamai.com/bot-manager/reference/put-javascript-injection-rules). Updates an existing JavaScript injection rule.

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

resource "akamai_botman_javascript_injection" "javascript_injection" {
  config_id            = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id   = "gms1_134637"
  javascript_injection = file("${path.module}/javascript_injection.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the JavaScript injection rule.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the JavaScript injection rule.
- `javascript_injection` (Required). JSON-formatted collection of JavaScript injection settings and setting values. In the preceding sample code, the syntax `file("${path.module}/javascript_injection.json")` points to the location of a JSON file containing the JavaScript injection settings and values.

**See also**

- [Set up JavaScript injection](https://techdocs.akamai.com/bot-manager/docs/set-up-javascript-injection)
