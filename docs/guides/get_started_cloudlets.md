---
layout: "akamai"
page_title: "Module: Cloudlets"
description: |-
  Cloudlets module for the Akamai Terraform Provider
---

# Cloudlets Guide

A Cloudlet is a value-added application that complements Akamai's core delivery solutions to solve specific business challenges. Cloudlets bring a site's business logic closer to the end user by placing it on the edge of the content delivery platform. Currently, this module supports two types of Cloudlets: Application Load Balancer and Edge Redirector.

For each Cloudlet instance on your contract, you can create any number of policies. Each policy is tied to a specific property. Policies are versioned. Within a policy version you define rules. These rules determine what happens when an incoming request matches on the criteria for triggering the Cloudlet.

Also, each Cloudlet has its own special behavior in Property Manager. A behavior contains the settings for handling and processing the content passing through the Akamai network. Behaviors are essentially the building blocks for a property's rules. A Cloudlet can start processing end-user requests only if the referenced properties have the appropriate behavior enabled in their rule tree and are active on the production network.

## Prerequisites

To activate a Cloudlet policy, you need to have at least one property created within the group and contract you manage your web assets with.

Use the [`akamai_property`](../resources/property.md) resource to create or import delivery properties.

### Property requirements for Cloudlets that forward requests

For these Cloudlets, you can set up a separate origin to forward incoming requests to:

* Application Load Balancer
* Audience Segmentation
* Forward Rewrite
* Phased Release

To add this type of origin to your property, you need to set up Conditional Origin rules and activate the property version containing these new rules. Once the activation is complete, then you can set up your Cloudlet-specific rules and behaviors. 

Follow these steps to add Conditional Origins to your Terraform configuration:

1. In the JSON rule tree file for the property that you want to use with your Cloudlet, set the mandatory [`origin` behavior](https://developer.akamai.com/api/core_features/property_manager/vlatest.html#origin)  to support Conditional Origins. 

    ~> **Note** Application Load Balancer uses Conditional Origins to represent data centers that are part of the load balancing and failover scheme. You should have a separate Conditional Origin for each data center included in your Application Load Balancing configuration.

2. Use the [`akamai_property_activation`](../resources/property_activation.md) resource to activate the property version containing the new Conditional Origins on the production network.

## Cloudlets workflow

With the Cloudlets module, you can create and activate new policies for your Cloudlets. All Cloudlets except for Application Load Balancer follow the same workflow. Application Load Balancer has additional steps for the load balancing settings.

### Basic workflow for Cloudlet policies

Use this workflow for all Cloudlets except [Application Load Balancer](#workflow-for-an-application-load-balancer-cloudlet-policy):

* [Optionally, configure the policy match rule with a dedicated data source](#configure-the-match-rule-with-a-data-source). This generates a JSON-encoded set of conditions that you can reference later on.
* [Get the group ID](#get-the-group-id). You need your group ID to create a policy.
* [Import or create a policy](#import-or-create-a-policy). This serves as a container for the rules that govern how a Cloudlet behaves.
* [Activate a policy](#activate-a-policy-version). Test the changes on staging and eventually push them to the production network.
* [Update the associated property](#update-the-associated-property). Set up the behavior for each Cloudlet you're configuring.
* [Activate the property version](#activate-the-property-version). After you modify the rule tree, activate the changed property on the production network.

### Workflow for an Application Load Balancer Cloudlet policy

* [Import or create the Application Load Balancer configuration](#import-or-create-the-application-load-balancer-configuration). Set up the traffic management details.
* [Activate the configuration](#activate-the-application-load-balancer-configuration). Push a version of the load balancing configuration to the network you select.
* [Optionally, configure the policy match rule with a dedicated data source](#configure-the-match-rule-with-a-data-source). This generates a JSON-encoded set of conditions that you can reference later on.
* [Get the group ID](#get-the-group-id). You need your group ID to create a policy.
* [Import or create a policy](#import-or-create-a-policy). This serves as a container for the rules that govern how a Cloudlet behaves.
* [Activate a policy](#activate-a-policy-version). Test the changes on staging and eventually push them to the production network.
* [Update the associated property](#update-the-associated-property). Set up a behavior for each Cloudlet you're configuring.
* [Activate the property version](#activate-the-property-version). After you modify the rule tree, activate the changed property on the production network.

## Import or create the Application Load Balancer configuration

You can either import an existing configuration or create a new one with Terraform:

### Import a configuration

If you already have an existing load balancing configuration you'd like to work with, you can use one of these import options:
* Prepare the HCL configuration and run `terraform import`. See the [`akamai_cloudlets_application_load_balancer`](../resources/cloudlets_application_load_balancer.md) resource.
* Use the [Terraform CLI](https://github.com/akamai/cli-terraform) to export the existing infrastructure to HCL and run `terraform import`. See the [`akamai_cloudlets_application_load_balancer`](../resources/cloudlets_application_load_balancer.md) resource.

### Create a configuration

Use the [`akamai_cloudlets_application_load_balancer`](../resources/cloudlets_application_load_balancer.md) resource to set up a load balancing configuration that includes the data centers and liveness tests you want to use.

## Activate the Application Load Balancer configuration

Use the [`akamai_cloudlets_application_load_balancer_activation`](../resources/cloudlets_application_load_balancer.md) resource to deploy the load balancing configuration version to either the Akamai staging or production network. You can activate a specific version multiple times if you need to.

Before activating on production, activate on staging first. This way you can detect any problems in staging before your changes progress to production.

## Configure the match rule with a data source

Optionally, to simplify the process, you can use the [`akamai_cloudlets_edge_redirector_match_rule`](../data-sources/cloudlets_edge_redirector_match_rule.md) or [`akamai_cloudlets_application_load_balancer_match_rule`](../data-sources/cloudlets_application_load_balancer_match_rule.md) data source to specify the match rules that govern how the Cloudlet is used. The data source returns the JSON-encoded rules for the Cloudlet that you can reference while creating a policy.

## Get the group ID

When setting up Cloudlets, you need to get the Akamai [`group_id`](../data-sources/group.md).

-> **Note** The Cloudlets module supports both ID formats, either with or without the `grp_` prefix. For more information about prefixes, see the ID prefixes section of the [Property Manager API (PAPI)](https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes) documentation.

## Import or create a policy

You can either import an existing policy or create a new one with Terraform:

### Import a policy

If you already have an existing policy you'd like to work with, you can use one of these import options:
* Prepare the HCL configuration and run `terraform import`. See the [`akamai_cloudlets_policy` resource](../resources/cloudlets_policy.md).
* Use the [Terraform CLI](https://github.com/akamai/cli-terraform) to export the existing infrastructure to HCL and run `terraform import`. See the [`akamai_cloudlets_policy` resource](../resources/cloudlets_policy.md).

### Create a policy

Use the [`akamai_cloudlets_policy`](../resources/cloudlets_policy.md) resource to configure and version a new policy to meet your business need. Within a single policy, you can specify up to 5000 rules that you want to apply. Either specify them directly using the resource attributes, or reference the JSON-encoded output of the [`akamai_cloudlets_edge_redirector_match_rule`](../data-sources/cloudlets_edge_redirector_match_rule.md) or [`akamai_cloudlets_application_load_balancer_match_rule`](../data-sources/cloudlets_application_load_balancer_match_rule.md) data source.

## Activate a policy version

Once you're satisfied with a policy version, use the [`akamai_cloudlets_policy_activation`](../resources/cloudlets_policy_activation.md) resource. Here, you associate a Cloudlets policy version with one or more properties that are within a compatible group. An activation deploys the policy to either the Akamai Staging or Production network. You activate a specific version, but the same version can be activated separately more than once.

Run `terraform apply`. Terraform shows an overview of changes, so you can still go back and modify the configuration, or confirm to proceed. See [Command: apply](https://www.terraform.io/docs/commands/apply.html).

-> **Note** Once you activate a policy version on either of the networks, any change to the `akamai_cloudlets_policy` resource creates and activates a new version upon running `terraform apply`.

## Update the associated property

To implement an activated Cloudlet on the edge, you need to add the Cloudlet's behavior for your Cloudlet to the JSON rule tree file in each property associated with your policy. See the [Property Provisioning Module Guide](../guides/get_started_property.md) for detailed instructions.

If you wish to customize the settings, see the behavior for your Cloudlet in the Property Manager API (PAPI) documentation: 

| Cloudlet      | Behavior |
| ----------- | ----------- |
| API Prioritization | [`apiPrioritization`](https://techdocs.akamai.com/property-mgr/reference/latest-apiprioritization) |
| Application Load Balancer   | [`applicationLoadBalancer`](https://techdocs.akamai.com/property-mgr/reference/latest-applicationloadbalancer) |
| Audience Segmentation | [`audienceSegmentation`](https://techdocs.akamai.com/property-mgr/reference/latest-audiencesegmentation) |
| Edge Redirector | [`edgeRedirector`](https://techdocs.akamai.com/property-mgr/reference/latest-edgeredirector) |
| Forward Rewrite | [`forwardRewrite`](https://techdocs.akamai.com/property-mgr/reference/latest-forwardrewrite) |
| Phased Release | [`phasedRelease`](https://techdocs.akamai.com/property-mgr/reference/latest-phasedrelease) |
| Request Control | [`requestControl`](https://techdocs.akamai.com/property-mgr/reference/latest-requestcontrol) |
| Visitor Prioritization | [`visitorPrioritization`](https://techdocs.akamai.com/property-mgr/reference/latest-visitorprioritization) |

## Activate the property version

Use the [`akamai_property_activation`](../resources/property_activation.md) resource to activate the modified property version on the production network. An active Cloudlet policy won't be applied to the user requests unless you activate the property version with appropriate Cloudlet behavior enabled.

Run `terraform apply` again to implement the changes.
