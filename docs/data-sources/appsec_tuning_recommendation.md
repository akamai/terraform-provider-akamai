---
layout: "akamai"
page_title: "Akamai: TuningRecommendation"
subcategory: "Application Security"
description: |-
 TuningRecommendation
---

# akamai_appsec_tuning_recommendation

**Scopes**: Security policy; attack group

Returns tuning recommendations for the specified attack group (or, if the **attack_group** argument is not included, returns tuning recommendations for all the attack groups in the specified security policy). Tuning recommendations help minimize the number of false positives triggered by a security policy; with a false positive, a client request is marked as having violated the security policy restrictions even though it didn't actually violate those restrictions. Tuning recommendations are returned as attack group exceptions: if you choose, you can copy the response and use the [akamai_appsec_attack_group](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_attack_group) resource to add the recommended exception to a security policy or attack group

If the data source response is empty, that means that there are no further recommendations for tuning your security policy or attack group.

If you need to, you can manually merge a recommended exception for an attack group with the exception previously configured in the attack group resource. 
Additional information is available in the [Application Security API v1 documentation](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getrecommendations).

**Related API endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/recommendation](https://developer.akamai.com/api/cloud_security/application_security/v1.html#gettuningrecommendationsforanattackgroup)s

## Example usage

```
terraform {
 required_providers
  akamai  = {
   source = "akamai/akamai"

 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

// USE CASE: user wants to view tuning recommendations for the specified security policy or for the specified attack group

data "akamai_appsec_configuration" "configuration" {
 name = "Documentation"
}

data "akamai_appsec_tuning_recommendations" "policy_recommendations" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637"
}

output "policy_recommendations_json" {
 value = data.akamai_appsec_tuning_recommendations.policy_recommendations.jso
}

data "akamai_appsec_tuning_recommendations" "attack_group_recommendations" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637"
 attack_group       = "SQL"
}

output "attack_group_recommendations_json" {
 value = data.akamai_appsec_tuning_recommendations.attack_group_recommendations.json
}
```



## Argument reference

This data source supports the following arguments:

- **config_id** (required). Unique identifier of the security configuration you want to return tuning recommendations for.
- **security_policy_id** (required). Unique identifier of the security policy you want to return tuning recommendations for.
- **attack_group** (optional). Unique name of the attack group you want to return tuning recommendations for. If not included, recommendations are returned for all your attack groups.



## Output options

The following options can be used to determine the information returned, and how that returned information is formatted:

- **json**. JSON-formatted list of the tuning recommendations for the security policy or the attack group.
