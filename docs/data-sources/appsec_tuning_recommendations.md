---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_tuning_recommendations

Returns tuning recommendations for the specified attack group or rule (or, if both the `attack_group` and the `rule_id` arguments are not included, returns tuning recommendations for all the attack groups and rules in the specified security policy).
Tuning recommendations help minimize the number of false positives triggered by a security policy. With a false positive, a client request is marked as having violated the security policy restrictions even though it actually did not.
Tuning recommendations are returned as attack group or rule exceptions: if you choose, you can copy the response and use the `akamai_appsec_attack_group` resource to add the recommended exception to an attack group or the `akamai_appsec_rule` resource to add the recommended exception to a rule.  
If the data source response is empty, that means that there are no further recommendations for tuning your security policy or attack group.
If you need, you can manually merge a recommended exception for an attack group or a rule with the exception previously configured.
You can find additional information in our [Application Security API v1 documentation](https://techdocs.akamai.com/application-security/reference/get-recommendations).

**Related API endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/recommendation](https://techdocs.akamai.com/application-security/reference/get-recommendations)

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_tuning_recommendations" "policy_recommendations" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "policy_recommendations_json" {
  value = data.akamai_appsec_tuning_recommendations.policy_recommendations.json
}

data "akamai_appsec_tuning_recommendations" "attack_group_recommendations" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  ruleset_type       = var.ruleset_type
  attack_group       = var.attack_group
}

output "attack_group_recommendations_json" {
  value = data.akamai_appsec_tuning_recommendations.attack_group_recommendations.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required). Unique identifier of the security configuration you want tuning recommendations for.

* `security_policy_id` - (Required). Unique identifier of the security policy you want tuning recommendations for.

* `ruleset_type` - (Optional). Type of ruleset used by the security configuration you want tuning recommendations for. Supported values are `active` and `evaluation`. Defaults to `active`.

* `attack_group` - (Optional). Unique name of the attack group you want tuning recommendations for. If both `attack_group` and `rule_id` not included, recommendations are returned for all attack groups.

* `rule_id` - (Optional). Unique id of the rule you want tuning recommendations for. If both `attack_group` and `rule_id` not included, recommendations are returned for all attack groups.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `json` - JSON-formatted list of the tuning recommendations for the security policy, the attack group or the rule. The exception block format in a recommendation conforms to the exception block format used in `condition_exception` element of `attack_group` or ASE rule resource.