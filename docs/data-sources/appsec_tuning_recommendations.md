---
layout: "akamai"
page_title: "Akamai: Tuning Recommendations"
subcategory: "Application Security"
description: |-
 TuningRecommendations
---

# akamai_appsec_tuning_recommendations

Use the `akamai_appsec_tuning_recommendations` data source to retrieve tuning recommendations for all attack groups in a security policy or for a specific attack group. To accept a recommendation, the exception block in the recommendation 
for an attack group, could be passed as condition_exception argument value of the attack group resource. 
If needed, the Akamai terraform provider users would have to merge a recommended exception for an attack group with the exception previously configured in the attack group resource. 
Additional information is available [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getrecommendations).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view tuning recommendations for a security policy
// /appsec/v1/configs/{configId}/versions/{versionNum}/security-policies/{policyId}/recommendations
// user wants to view tuning recommendations for an attack group if attack group is specified
// /appsec/v1/configs/{configId}/versions/{versionNum}/security-policies/{policyId}/recommendations/attack-group/{attack-group-id} 
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

* `config_id` - (Required) The configuration ID.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - (Optional) The ID of the attack group to use. If not supplied, information about all attack groups will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted information about tuning recommendations. The exception block format in a recommendation conforms to the exception block format used in condition_exception element of attack_group resource.  


