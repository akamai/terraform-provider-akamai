---
layout: "akamai"
page_title: "Module: EdgeWorkers"
description: |-
  EdgeWorkers module for the Akamai Terraform Provider
---

# EdgeWorkers Guide

With this module you can create EdgeWorkers functions to run JavaScript at the edge of the Internet to dynamically manage web traffic. You can use the Akamai Terraform provider to deploy custom code on thousands of edge servers and apply logic that creates powerful web experiences.

## Prerequisites

Before you start, make sure Edgeworkers is in your contract. You can find the list of products within your account in ​Control Center​ under Contracts. Contact your ​Akamai​ support team to enable Edgeworkers if necessary.

Complete the tasks in [Get Started with the Akamai Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider). You should have an API client and a valid `akamai.tf` Terraform configuration before adding the EdgeWorkers module configuration.

## EdgeWorkers workflow

* [Get the group ID](#get-the-group-id). You need your group ID to create a EdgeWorker ID.
* [Get the resource tier](#get-the-resource-tier). Choose the resource tier you want to use for your EdgeWorker ID.
* [Add data to the EdgeKV namespace](#add-data-to-the-edgekv-namespace). You can optionally add an EdgeKV key-value store to your EdgeWorkers function.
* [Create an EdgeWorker ID](#create-an-edgeworker-id). Create a new EdgeWorker ID.
* [Add an EdgeWorker version](#add-an-edgeworker-version). Add a version to your EdgeWorker ID.
* [Create a property rule](#create-a-property-rule) Create a rule in your property.
* [Activate the EdgeWorker version](#activate-the-edgeworker-version). After you modify the rule tree, activate the changed property on the production network.
* [Test the EdgeWorker version](#test-the-edgeworker-version). Use the instructions in the the Akamai documentation to test your EdgeWorker version.

## Get the group ID

When creating an EdgeWorker ID, you need to get the Akamai [`group_id`](../data-sources/group.md).

-> **Note** The EdgeWorkers module supports both ID formats, either with or without the `grp_` prefix. For more information about prefixes, see the [ID prefixes](https://techdocs.akamai.com/property-mgr/reference/id-prefixes) section of the Property Manager API (PAPI) documentation.

## Get the resource tier

Use the [akamai_edgeworkers_resource_tier](../data-sources/edgeworkers_resource_tier.md) data source to get detailed information about resource tiers available for your contract. Store the `resource_tier_id` you want to use in your EdgeWorker ID.

-> **Note** To include an EdgeKV database in your EdgeWorkers function you need to create an EdgeWorker ID using the Dynamic Compute resource tier. To do this set the `resource_tier_id` to 200 when you create your EdgeWorker ID.

## Create an EdgeWorker ID

Create an EdgeWorker ID that includes a name, a group id, and a resource tier.

Use the [`akamai_edgeworker`](../resources/edgeworker.md) resource to set up an EdgeWorker ID you need to create and activate EdgeWorker versions.

## Add data to the EdgeKV namespace

You can optionally use the [`akamai_edgekv`](../resources/edgekv.md) resource to add an EdgeKV database to your EdgeWorkers function.

For instructions on how to create an EdgeKV-enabled EdgeWorkers code bundle refer to the [EdgeKV documentation](https://techdocs.akamai.com/edgekv/docs/create-a-code-bundle).

-> **Note** Customers are responsible for maintaining control over the data hosted on this service and for appropriately using the data returned by EdgeKV. EdgeKV does not support storage of sensitive information where the consequence of an unauthorized disclosure would be a serious business or compliance issue.

Customers should not use sensitive information when creating namespaces, groups, keys, or values. EdgeKV currently encrypts data at rest and in transit but does not protect against unauthorized disclosure. The requisite processes and controls to enable the storage of sensitive data will be available in a future release.

## Add an EdgeWorker version

The `terraform apply` command automatically creates the first EdgeWorker version. Each time you modify the EdgeWorkers code bundle and use the `terraform apply` command a new EdgeWorker version is created.

For instructions on how to create a code bundle refer to the [EdgeWorkers documentation](https://techdocs.akamai.com/edgeworkers/docs/create-a-code-bundle).

## Activate the EdgeWorker version

Use the [`akamai_edgeworkers_activation`](../resources/edgeworkers_activation.md) resource to activate the property containing your EdgeWorker version. This operation takes approximately 20 minutes.

## Create a property rule

Use the [`akamai_edgeworkers_property_rules`](../data-sources/edgeworkers_property_rules.md) to create a new property rule.

## Activate the property version

To activate an EdgeWorker version you need to add the EdgeWorkers behavior to the JSON rule tree file in each property associated with your policy. See the [Property Provisioning Module Guide](../guides/get_started_property.md) for detailed instructions.

If you wish to customize the settings, see the [`edgeWorker` behavior](https://techdocs.akamai.com/property-mgr/reference/latest-edgeworker) in the PAPI Catalog Reference.

Use the [`akamai_edgeworkers_activation`](../resources/edgeworkers_activation.md) resource to activate the modified property version on the staging or production network.

Run `terraform apply` again to implement the changes.

## Test the EdgeWorker version

Once your property is active use the instructions in the [EdgeWorkers documentation](https://techdocs.akamai.com/edgeworkers/docs/test-hello-world-2) to test your EdgeWorker version.
