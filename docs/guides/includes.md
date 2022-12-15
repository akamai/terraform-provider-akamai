---
layout: "akamai"
page_title: "Property Includes"
description: |-
   Includes integration for Property Provisioning for the Akamai Terraform Provider
---

# Includes

Includes are small chunks of a property configuration that you can create, version, and activate independently from the rest of the property's rule tree. With includes, you can delegate different parts of single domains to responsible teams or implement common settings you can share across multiple properties.

## Before you begin

* Understand the [basics of Terraform](https://learn.hashicorp.com/terraform?utm_source=terraform_io).
* Complete the steps in [Get started](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started) and [Set up your authentication](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/auth).
* Set up your [property and its configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property).

## How it works

![Workflow for includes: create, activate, add](https://techdocs.akamai.com/terraform-images/img/includes.svg)

## Create include

You can create a new include from scratch or base one on an existing include version. Each include you create is one of two types, `MICROSERVICES` or `COMMON_SETTINGS`.

|Type|Description|Example|
|---|---|---|
|`MICROSERVICES`|Use when different teams will work independently on different parts of a single site. For each include of this type:<ul><li>Your users can set up automation for the include and test and deploy on their own schedule.</li><li>You can control access to an include with [groups](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_iam#create-groups), so that an application team can only edit the settings within a designated include, without having access to the parent property.</li></ul>You can create up to 300 of this type per parent property.|Teams managing and deploying their own set of rules.|
|`COMMON_SETTINGS`|Use when a central team manages property configurations for your site. This type of include functions as a single source of information for common settings used across different configurations.<br />You can create up to three of this type per parent property.|Instead of repeating work for a number of settings used by both site A and site B, changes to an include managing the shared settings updates both sites.|

The rules in your includes also follow a versioned rule format. You should work with the [**most recent frozen rule format**](https://techdocs.akamai.com/property-mgr/reference/get-rule-formats) as includes do not work with the **latest** set of rule formats.

<blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">Includes and parent properties must use the same rule format.</div></blockquote>

### New with specific rule tree

You can create a new include with a specific rule tree. This rule tree can be a new one you've configured or one from an existing include that you've edited.

1. Create or edit a rule tree.

    * To create a new rule tree, use the [akamai_property_rules_template](../data-sources/property_rules_template.md) data source.
    * To use and edit an existing include's rule tree:

        1. Get an existing include's rules with the [akamai_property_include_rules](../data-sources/property_include_rules.md) data source.
        1. Create a new JSON file with the returned rules and add your edits.
        1. Use your new JSON as a template file in the [akamai_property_rules_template](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/property_rules_template#how-to-work-with-json-template-files) data source.

1. Once you've got a rule tree, use the values from the `akamai_property_rules_template` as variables in the [akamai_property_include](../resources/property_include.md) resource to create your new include. The `akamai_property_rules_template` data source lets you define a rule tree.

    ```hcl
    resource "akamai_property_include" "new_specific_rule_tree" {
        contract_id = "C-0N7RAC7"
        group_id    = "X112233"
        product_id  = "prd_123456"
        name        = "example3"
        rule_format = "v2022-10-18"
        rules       = data.akamai_property_rules_template.example.json // Here "example" would be the local name of your data source.
    }
    ```

1. Run `terraform plan` to check your syntax and review your changes and then run `terraform apply` to create an unactivated, editable include.
1. When you're ready, [activate](#activate-include) your include.

### New based on existing

Each iteration of an include increases its version number, and all of those versions are available to use as a starting point for a new include, no matter its activation status. Creating a new include this way uses an existing versions rules. You can use the rule set as is or use it as your base for an iteration.

<center><img src="https://techdocs.akamai.com/terraform-images/img/new-from-previous.png" height="150px"; width="450px"; alt="You can create a new include version using any of its previous versions as a starting point." /></center>

1. Using the [akamai_property_includes](../data-sources/property_includes.md) data source, run `terraform plan` to get all of the includes for your contract and group. This will give the include IDs and versions. Even though the parent's and include's rule formats needs to match, if the rule format of parent is the latest, this data source will list all the includes under the contract.

    ```hcl
    data "akamai_property_includes" "get_all_includes" {
        contract_id = "C-0N7RAC7"
        group_id    = "X112233"
    }

    output "get_all_includes" {
        value = data.akamai_property_includes.get_all_includes
    }
    ```

1. Add the ID and version of the include you want to use to your group and contract IDs in the [akamai_property_include_rules](../data-sources/property_include_rules.md) data source to identify the rule set and format you want to use.

    ```hcl
    data "akamai_property_include_rules" "rules_v2" {
        include_id  = "inc_X12345"
        contract_id = "C-0N7RAC7"
        group_id    = "X112233"
        version     = 2
    }

    output "rules_v2" {
        value = data.akamai_property_include_rules.rules_v2
    }
    ```

1. If you want to use the rules as is, pass your contract and group IDs, a unique name, and variables for both the `rule_format` and `rules` in the [akamai_property_include](../resources/property_include.md) resource. If you have edited the rules, add in a pointer to the JSON file location instead and use the most recent version.

   ```hcl
    resource "akamai_property_include" "new_from_existing" {
        contract_id = "C-0N7RAC7"
        group_id    = "X112233"
        product_id  = "prd_123456"
        name        = "example2"
        rule_format = data.akamai_property_include_rules.rules_v2.rule_format
        rules       = data.akamai_property_include_rules.rules_v2.rules
    }
    ```

1. Run `terraform plan` to check your syntax and review your changes and then run `terraform apply` to create an unactivated, editable include.
1. When you're ready, [activate](#activate-include) your include.

## Activate include

Because includes can be deployed independently of and directly affect your larger property configurations, the importance of checking their behavior on staging before deployment to your production instances cannot be overstated. For this reason, each of your includes should have one activation in staging and one in production.

* For *new* includes, activation does not immediately apply their rules because they're not directly assigned to hosts.
    Instead, activation of *new* include makes the include available for use.
    The actual application of their rules comes when you associate an includes to a property.

* For *existing* includes that are versioning up, activation immediately deploys the changes to the parent property.

To activate an include, send the include's version information and the network setting through the [akamai_property_include_activation](../resources/property_include_activation.md) resource.


1. Using the [akamai_property_includes](../data-sources/property_includes.md) data source, run `terraform plan` to get all of the includes for your contract and group. This will give the include IDs and versions.

    ```hcl
    data "akamai_property_includes" "get_all_includes" {
        contract_id = "C-0N7RAC7"
        group_id    = "X112233"
    }

    output "get_all_includes" {
        value = data.akamai_property_includes.get_all_includes
    }
    ```

1. Send the returned `include_id` for the one you want to activate as part of the required information in the `akamai_property_include_activation` resource.

    <blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">Compliance information is only required for an activation in a production.</div></blockquote>

    ```hcl
    resource "akamai_property_include_activation" "activate_include" {
      include_id    = "inc_X12345"
      contract_id   = "C-0N7RAC7"
      group_id      = "X112233"
      version       = 1
      network       = "STAGING"
      notify_emails = [
          "example@example.com",
          "example2@example.com"
      ]
      note = "An optional statment about your activation"
      compliance_record { // Required when activating in production only.
        noncompliance_reason       = "NONE"
        ticket_id                  = "JIRA-1234"
        customer_email             = "example@example.com"
        peer_reviewed_by           = "John Doe"
        unit_tested                = true
      }
    }
    ```

1. Run `terraform plan` to check your syntax and review your changes and then run `terraform apply` to activate your include. You can now add a new include to a property. An immediate deploy of changes for existing includes occurs in the designated environment.

## Add include to property

Your include's rules do not affect the behavior of a host until it's added to a property. Includes are connected to a property in a parent-child relationship.

* You can add up to 300 `MICROSERVICES` and three `COMMON_SETTINGS` includes to a single parent property.
* You can add an include to several different properties.

To add your include to a property, use the [property_rules_template](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/property_rules_template).

## Deactivate, reactivate, and delete

When you activate a new version of an include, it automatically deactivates the previous version. If you decide you want to use a previous version again, pass an include ID and a specfic version in the `akamai_property_include_activation` resource to reactivate it.

To delete an include, remove the include from all of its parent properties, deactivate the active versions of your include, and then use `terraform destroy` to delete it.

1. Use the [akamai_property_include_parents](../data-sources/property_include_parents.md) data source to find out which properties use the include.

     ```hcl
    resource "akamai_property_include" "get_all_parents" {
        group_id    = "X112233"
        contract_id = "C-0N7RAC7"
        include_id  = "inc_X12345"
    }
    ```

1. Edit the rule tree for each property returned and remove the includes from the rule JSON using the [property_rules_template](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/property_rules_template).
1. Run `terraform destroy` on the [akamai_property_include_activation](../resources/property_include_activation.md) resource.
1. Run `terraform destroy` on the [akamai_property_include](../resources/property_include.md) resource.
