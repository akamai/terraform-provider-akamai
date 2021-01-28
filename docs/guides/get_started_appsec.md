---
layout: "akamai"
page_title: "Akamai: Get Started with Application Security"
description: |-
   Application Security in Akamai provider for Terraform
---

# Application Security in Akamai provider for Terraform 

Application Security (appsec) in the Akamai Terraform provider (provider) enables application 
security configurations including the following: 
* custom rules.
* match targets.
* other application security resources that operate within the Cloud.

This Guide is for developers who:
* are interested in implementing or updating an integration of Akamai functionality with Terraform.
* already have some familiarity with Akamai.
* understand how to create and edit the 'akamai.tf' file [see](get_started_akamai.md).
 
For details about Akamai's application security, see the [API documentation](https://developer.akamai.com/api/cloud_security/application_security/v1.html)

## Prerequisites

To manage Application Security resources, you need to obtain information regarding your 
existing security implementation, including the following information:

* **Configuration ID**: The ID of the specific security configuration under which the resources are defined.

In many cases, you need additional information, which often includes the 
version number of the security configuration (see below).

### Retrieve existing security configuration information

You can obtain the name and ID of the existing security configurations by using the 
[`akamai_appsec_configuration`](../data-sources/appsec_configuration.md) data source. 
Using it without parameters outputs information about all security configurations associated with your account. 

Add the following to your `akamai.tf` file:

```hcl
data "akamai_appsec_configuration" "configurations" {
}

output "configuration_list" {
  value = data.akamai_appsec_configuration.configurations.output_text
}
```

Save the resulting text file, and then use terminal to initialize Terraform with the command:

```bash
$ terraform init
```

This installs the latest version of the Akamai provider, along with any other providers necessary. 

When you need to obtain an update of Akamai provider, run `terraform init` again.

## Configure the Provider

Set up your .edgerc credential files as described in [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started), and include read-write permissions for the Application Security API. 

1. Create a new folder called `terraform`
1. Inside the new folder, create a new file called `akamai.tf`.
1. Add the provider configuration to your `akamai.tf` file:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
  config_section = "appsec"
}
```

## Test Your Configuration

To test your configuration, use `terraform plan`:

```bash
$ terraform plan
```

This command causes Terraform to create a plan for the work it will do, based on the configuration file. This does *not* actually make any changes and is safe to run as many times as you like.

## Apply changes

To display existing configuration information, or to create or modify resources as described in this guide, tell Terraform to `apply` the changes outlined in the plan by running the command:

```bash
$ terraform apply
```

Terraform responds with a formatted list of all existing security configurations in your account, along with names and IDs (`config_id`), the most recently created version, and the version currently active in staging and production, if applicable.

When you have identified the desired security configuration by name, you can load that specific configuration into Terraform's state. 

To load a specific configuration:
1. Identify the desired security configuration by name, 
1. Edit your `akamai.tf` file to add the desired `name` parameter to the `akamai_appsec_configuration` data block.
1. Change the `output` block so that it gives just the `config_id` attribute of the configuration. 

After these changes, the section of your file below the initial `provider` block looks like the following example:

```hcl
data "akamai_appsec_configuration" "configuration" {
  name = "Example"
}

output "ID" {
  value = data.akamai_appsec_configuration.configuration.config_id
}
```

After running `terraform apply` on this file, the terminal displays `config_id` with the configuration value.

## Specify configuration to display

The provider's [`akamai_appsec_export_configuration`](../data-sources/appsec_export_configuration.md) data source can display complete information about any configuration that you specify, including attributes like custom rules, and selected hostnames. 

To show custom rule and selected hostname data for your most recent configuration, add the following blocks to your `akamai.tf` file:

```hcl
data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  search = [
  "customRules",
  "selectedHosts"
  ]
}

output "exported_configuration_text" {
  value = data.akamai_appsec_export_configuration.export.output_text
}
```

NOTE: You can specify any available version of the configuration. 
See the [`akamai_configuration_version`](../data-sources/appsec_configuration_version.md) 
data source to list the available versions. You can also specify other kinds of data for export 
using any of the following search fields:

* customRules
* matchTargets
* ratePolicies
* reputationProfiles
* rulesets
* securityPolicies
* selectableHosts
* selectedHosts

Save the file and run `terraform apply` to see a formatted display of the selected data.

## Add a hostname to the `selectedHosts` list

You can modify the list of hosts protected by a specific security configuration using 
the [`akamai_appsec_selected_hostnames`](../data-sources/appsec_selected_hostnames.md) resource. 
Add the following resource block to your `akamai.tf` file, replacing `example.com` with a hostname 
from the list reported in the `data_akamai_appsec_export_configuration` data source example above:

```hcl
resource "akamai_appsec_selected_hostnames" "selected_hostnames_append" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  hostnames = [ "example.com" ]
  mode = "APPEND"
}

output "selected_hostnames_appended" {
  value = akamai_appsec_selected_hostnames.selected_hostnames_append.hostnames
}
```

When you save the file and run `terraform apply`, Terraform updates the list of selected hosts and outputs the new list as values for `selected_hostnames_appended`. 

NOTE: You cannot modify a security configuration version that is currently active in staging or production, so the resource block above must specify an inactive version. 

After completing your changes to a security configuration version, you can activate it in staging.

## Activate a configuration version

To activate a specific configuration version, use the [`akamai_appsec_activations`](../resources/appsec_activations.md) resource. 

Add the following resource block to your `akamai.tf` file, replacing the `version` value with the number of a currently inactive version, perhaps the one you modified using the `akamai_appsec_selected_hostnames` resource above.

```hcl
resource "akamai_appsec_activations" "activation" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  network = "STAGING"
  notes  = "TEST Notes"
  notification_emails = [ "my_name@mycompany.com" ]
}
```

After you save the file and run `terraform apply`, Terraform activates the configuration version in staging. Upon completion of the activation, emails are sent to the addresses specified in the `notification_emails` list.


## Beta features

NOTE: The following data sources and resources are currently in Beta, and their behavior or documentation might change in a future release:

### Data Sources
  * akamai_appsec_eval
  * akamai_appsec_eval_rule_actions
  * akamai_appsec_eval_rule_condition_exception
  * akamai_appsec_ip_geo
  * akamai_appsec_rule_actions
  * akamai_appsec_rule_condition_exception
  * akamai_appsec_penalty_box
  * akamai_appsec_security_policy_protections
  * akamai_appsec_rate_policies
  * akamai_appsec_rate_policy_actions
  * akamai_appsec_rate_protections
  * akamai_appsec_reputation_protections
  * akamai_appsec_reputation_profiles
  * akamai_appsec_reputation_profile_actions
  * akamai_appsec_rule_upgrade_details
  * akamai_appsec_slow_post
  * akamai_appsec_slowpost_protections
  * akamai_appsec_attack_group_actions
  * akamai_appsec_waf_mode
  * akamai_appsec_waf_protection
  * akamai_appsec_attack_group_condition_exception

### Resources
  * akamai_appsec_eval
  * akamai_appsec_eval_rule_action
  * akamai_appsec_eval_rule_condition_exception
  * akamai_appsec_ip_geo
  * akamai_appsec_rule_condition_exception
  * akamai_appsec_rule_action
  * akamai_appsec_penalty_box
  * akamai_appsec_security_policy_protections
  * akamai_appsec_rate_policy
  * akamai_appsec_rate_policy_action
  * akamai_appsec_rate_protection
  * akamai_appsec_reputation_protection
  * akamai_appsec_reputation_profile
  * akamai_appsec_reputation_profile_action
  * akamai_appsec_rule_upgrade
  * akamai_appsec_slow_post
  * akamai_appsec_slowpost_protection
  * akamai_appsec_attack_group_action
  * akamai_appsec_waf_mode
  * akamai_appsec_waf_protection
  * akamai_appsec_attack_group_condition_exception
