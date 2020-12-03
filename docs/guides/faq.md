---
layout: "akamai"
page_title: "Akamai: FAQ (Frequently Asked Questions)"
description: |-
  Frequently Asked Questions
---

# Frequently Asked Questions

## Primary Zone Partially Created

In the rare instance, a primary zone may be only partially created on the Akamai backend; for example in the case of a network error. In this situation, the zone may have been created but not the SOA and NS records. Hence forth, any attempt to manage or administer recordsets in the zone will fail. The SOA and NS records must be manually created in order to continue to manage the configuration.

The records can be created either thru the [Akamai Control Center](https://control.akamai.com) or via the [CLI-DNS](https://github.com/akamai/cli-dns) package for the [Akamai CLI](https://developer.akamai.com/cli).
 
## Migrating an Edge DNS Zone and Records to Terraform

Migrating an existing Edge DNS Zone can be done in many ways. Two such methods include:

### Via Command Line Utility

A package, [CLI-Terraform](https://github.com/akamai/cli-terraform), for the [Akamai CLI](https://developer.akamai.com/cli) provides a time saving means to collect information about, generate a configuration for, and import an existing Edge DNS Zone and its contained recordsets. With the package, you can:

1. Generate a json formatted list of the zone and recordsets
2. Generate a Terraform configuration for the zone and select recordsets
3. Generate a command line script to import all defined resources

#### Notes
1. Terraform limits the characters that can be part of it's resource names. During construction of the resource configurations invalid characters are replaced with underscore , '_'
2. Terrform does not have any state during import of resources. Discrepencies may be identified in certain field lists during the first plan and/or apply following import as Terraform reconciles configurations and state. These discrepencies will clear following the first apply. 
3. The first time plan or apply is run, an update will be shown for the provider defined zone fields: contract and group.

It is recommended that the existing zone configuration and master file (using the API or Control Center) be downloaded before hand as a backup and reference.  Additionally, a terraform plan should be executed after importing to validate the generated tfstate.

### Via Step By Step Construction

1. Download your existing zone master file configuration (using the API) as a backup and reference.
2. Using the zone master file as a reference, create a Terraform configuration representing the the existing zone and all contained recordsets. Note: In creating each resource block, make note of `required`, `optional` and `computed` fields.
3. Use the Terraform Import command to import the existing zone and contained recordsets; singularly and in serial order.
4. (Optional, Recommended) Review and compare the zone master file content and created Terraform.tfstate to confirm the zone and all recordsets are represented correctly
5. Execute a `Terraform Plan` on the configuration. The plan should be empty. If not, correct accordingly and repeat until plan is empty and configuration is in sync with the Edge DNS Backend.

Since Terraform assumes it is the de-facto state for any resource it leverages, we strongly recommend staging the zone and recordset imports in a test environment to familiarize yourself with the provider operation and mitigate any risks to the existing Edge DNS zone configuration.

## Migrating a GTM domain (and contained objects) to Terraform

Migrating an existing GTM domain can be done in many ways. Two such methods include:

### Via Command Line Utility

A package, [CLI-Terraform](https://github.com/akamai/cli-terraform), for the [Akamai CLI](https://developer.akamai.com/cli) provides a time saving means to collect information about, generate a configuration for, and import an existing GTM domain and its contained objects and attributes. With the package, you can:

1. Generate a json formatted list of all domain objects
2. Generate a Terraform configuration for the domain and contained objects
3. Generate a command line script to import all defined resources

#### Notes
1. Terraform limits the characters that can be part of it's resource names. During construction of the resource configurations invalid characters are replaced with underscore , '_'
2. Terrform does not have any state during import of resources. Discrepencies may be identified in certain field lists during the first plan and/or apply following import as Terraform reconciles configurations and state. These discrepencies will clear following the first apply. 
3. The first time plan or apply is run, an update will be shown for the provider defined domain fields: contract, group and wait_on_complete.

It is recommended that the existing domain configuration (using the API or Control Center) be downloaded before hand as a backup and reference.  Additionally, a terraform plan should be executed after importing to validate the generated tfstate.

### Via Step By Step Construction

1. Download your existing domain configuration (using the API or Control Center) as a backup and reference.
2. Using the domain download as a reference, create a Terraform configuration representing the existing domain and all contained GTM objects. Note: In creating each resource block, make note of `required`, `optional`, and `computed` fields.
3. Run `terraform import`. This command imports the existing domain and contained objects one at a time based on the order in the configuration.
4. (Optional, Recommended) Review domain download content and created terraform.tfstate to confirm the domain and all objects are represented correctly
5. Run `terraform plan` on the configuration. The plan should be empty. If not, correct accordingly and repeat until plan is empty and configuration is in sync with the GTM Backend.

Since Terraform assumes it is the de facto state for any resource it leverages, we strongly recommend staging the domain and objects imports in a test environment to familiarize yourself with the provider operation and mitigate any risks to the existing GTM domain configuration.

## GTM Module: Resource fields when using plan or apply commands

When using `terraform plan` or `terraform apply`, Terraform presents both fields defined in the configuration and all defined resource fields. Fields are either required, optional or computed as specified in each resource description. Default values for fields will display if not explicitly configured. In many cases, the default will be zero, empty string, or empty list depending on the the type. These default or empty values are informational and not included in resource updates.
