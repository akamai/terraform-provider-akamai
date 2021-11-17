---
layout: "akamai"
page_title: "Akamai: Tuning Recommendations"
subcategory: "Application Security"
description: |-
 TuningRecommendations
---

# akamai_appsec_tuning_recommendations

Returns tuning recommendations for the specified attack group (or, if the `attack_group` argument is not included, returns tuning recommendations for all the attack groups in the specified security policy).
Tuning recommendations help minimize the number of false positives triggered by a security policy. With a false positive, a client request is marked as having violated the security policy restrictions even though it actually did not.
Tuning recommendations are returned as attack group exceptions: if you choose, you can copy the response and use the `akamai_appsec_attack_group` resource to add the recommended exception to a security policy or attack group.
If the data source response is empty, that means that there are no further recommendations for tuning your security policy or attack group.
If you need, you can manually merge a recommended exception for an attack group with the exception previously configured in the attack group resource. 
You can find additional information in our [Application Security API v1 documentation](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getrecommendations).

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
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "policy_recommendations_json" {
  value = data.akamai_appsec_tuning_recommendations.policy_recommendations.json
}

data "akamai_appsec_tuning_recommendations" "attack_group_recommendations" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
}

output "attack_group_recommendations_json" {
  value = data.akamai_appsec_tuning_recommendations.attack_group_recommendations.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required). Unique identifier of the security configuration you want to return tuning recommendations for.

* `security_policy_id` - (Required). Unique identifier of the security policy you want to return tuning recommendations for.

* `attack_group` - (Optional). Unique name of the attack group you want to return tuning recommendations for. If not included, recommendations are returned for all your attack groups.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `json` - JSON-formatted list of the tuning recommendations for the security policy or the attack group. The exception block format in a recommendation conforms to the exception block format used in `condition_exception` element of `attack_group` resource.  


