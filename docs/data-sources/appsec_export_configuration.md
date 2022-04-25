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

The `search` parameter enables you to limit your search to specific resources. For example, to search only for attack groups, include the search parameter in your Terraform configuration:

```
search = ["attackGroups"]
```

In turn, Terraform exports information about all your attack groups (note that the following output has been truncated to save space):

```
+----------------------------------------------------------------------------+
| attackGroups                                                               |
+----------+---------------------------------+----------+--------------------+
| ID       | NAME                            | TYPE     | RULESET VERSION ID |
+----------+---------------------------------+----------+--------------------+
| CMD      | Command Injection               | ASE AUTO | 7257               |
| LFI      | Local File Inclusion            | ASE AUTO | 7257               |
| OUTBOUND | Total Outbound                  | ASE AUTO | 7257               |
| PLATFORM | Web Platform Attack             | ASE AUTO | 7257               |
| POLICY   | Web Policy Violation            | ASE AUTO | 7257               |
```

This type of output is produced by setting the value of the `search` parameter to the appropriate search term (in this case, **attackGroups**).

However, it’s also possible for the akamai_appsec_export_configuration data source to return data formatted in a way that makes it easy to import that data back into Terraform. With this approach, the output for a single attack group might look similar to this (note the presence of the Terraform `import` command):

```
// terraform import akamai_appsec_attack_group.akamai_appsec_attack_group_gms1_134637 51219:gms1_134637:SQL

resource "akamai_appsec_attack_group" "akamai_appsec_attack_group_gms1_134637" {
          config_id = 51219
          security_policy_id = "gms1_134637"
          attack_group = "SQL"
          attack_group_action = "alert"
    }
```

To get this type of output, add **.tf** to end of your search term. For example, instead of searching for attackGroups, add **.tf** and search for **AttackGroups.tf**:

```
search = ["AttackGroup.tf"]
```

Note that even though you seem to be specifying the name of a Terraform configuration file (e.g., `activeGroups.tf`) this data source does not save the exported data to a text file. In other words, running the command shown above won’t create a file named `activeGroups.tf` that contains all the exported data.


**Related API Endpoint**: [/appsec/v1/export/configs/{configId}/versions/{versionNumber}](https://techdocs.akamai.com/application-security/reference/get-export-config-version)

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
- `search` (Optional). JSON array of strings specifying the types of information to be retrieved. Note that there are two different ways to return data by using the `search` parameter. To return data in tabular format, use one or more of the following terms:

   - attackGroups
   - customDenyList
   - customRules
   - matchTargets
   - ratePolicies
   - reputationProfiles
   - rules
   - securityPolicies
   - selectableHosts
   - selectedHosts

To return data that can be easily imported back into Terraform, use one or more of these terms:

    - AdvancedSettingsLogging.tf
    - AdvancedSettingsEvasivePathMatch.tf
    - AdvancedSettingsPragmaHeader.tf
    - AdvancedSettingsPrefetch.tf
    - ApiRequestConstraints.tf
    - CustomDeny.tf
    - CustomRule.tf
    - CustomRuleAction.tf
    - MatchTarget.tf
    - PenaltyBox.tf
    - RatePolicy.tf
    - RatePolicyAction.tf
    - ReputationProfile.tf
    - ReputationProfileAction.tf
    - Rule.tf
    - EvalRule.tf
    - AttackGroup.tf
    - EvalGroup.tf
    - ThreatIntel.tf
    - SecurityPolicy.tf
    - SelectedHostname.tf
    - SiemSettings.tf
    - SlowPost.tf
    - IPGeoFirewall.tf
    - WAPSelectedHostnames.tf


## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. Complete set of information about the specified security configuration version in JSON format. When this option is included information is always returned for the _entire_ configuration. Among other things, that means that, if your command uses the `search` parameter, that parameter is ignored.
- `output_text`. Tabular report showing the types of data specified in the `search` parameter. Valid only if the `search` parameter references at least one type.