---
layout: "akamai"
page_title: "Module: Image and Video Manager"
description: |-
  Image and Video Manager module for the Akamai Terraform Provider
---

# Image and Video Manager Guide (Beta)

The Image and Video Manager module lets you manage, optimize, and modify your images and short-form videos.
With Image and Video Manager, you can deliver a smooth, efficient, visually engaging customer experience that can also save you time and money.

~> Note: This functionality is currently in beta.

## About policies and policy sets

Image and Video Manager policies contain the settings to optimize and transform your images or videos.

A policy set is a collection of policies. An image policy set can contain only image policies, and a video policy set can contain only video policies.

~> Note: For more information about policies and policy sets, see [Create and edit policies](https://techdocs.akamai.com/ivm/docs/create-edit-policies) in the Image and Video Manager documentation.

### Policy and policy set relationships

Some things to keep in mind about policy sets and policies:

* A policy set belongs to a single contract.
* A policy belongs to a single policy set.
* A policy set can be used across multiple properties.
* A property can have multiple policy sets.

### The .auto policy

When you create a policy set, it automatically creates a default `.auto` policy.
This policy provides the baseline settings that
determine how the policy generates derivative images and videos.

### Activation considerations

By default, when you save a policy set, it's automatically activated on both staging and production networks.
Nothing happens to the policy set until you link it to a property.

When you save a policy, it's automatically activated on staging.
To activate a policy on production, you have to set the `activate_on_production` flag to `true` in either the [imaging_policy_image](../resources/imaging_policy_image.md) or [akamai_imaging_policy_video](../resources/imaging_policy_video.md) resource, and save your change.

### Variables

Many Image and Video Manager arguments let you specify a variable object instead of a string, number, or boolean value.

When using variables, you define the variable name in an argument that ends in `_var`. For example, if you want to have a variable for the gravity setting in a transformation, you’d use the `gravity_var` argument, not the `gravity` one.

## Prerequisites

Before you start, make sure Image and Video Manager is in your contract, and your contract includes the type of media (images or videos or both) that you intend to work with.
You can find the list of products within your account in ​Control Center​ under Contracts.
Contact your ​Akamai​ support team to enable Image and Video Manager if necessary.

Complete the tasks in [Get Started with the Akamai Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider).
You should have an API client and a valid Terraform configuration before adding this module.

## Image and Video Manager workflow

* [Get the contract ID](#get-the-contract-id)
* [Export a policy set and related policies](#export-a-policy-set-and-related-policies)
* [Add policy sets](#add-policy-sets)
* [Add policies](#add-policies)
* [Update the property rule tree](#update-the-property-rule-tree)
* [Activate the property on staging](#activate-the-property-on-staging)
* [Test your images and videos](#test-your-images-and-videos)
* [Activate the policy on production](#activate-the-policy-on-production)
* [Activate the property on production](#activate-the-property-on-production)

## Get the contract ID

When setting up Image and Video Manager, you need to retrieve the Akamai ID for your contract.
You can use the [akamai_contract](../data-sources/contract.md) data source to get this ID.

~> **Note:** If you use prefixes with your IDs, you might have to remove the `ctr_` prefix from your entry.
For more information about prefixes, see the [ID prefixes](https://techdocs.akamai.com/property-mgr/reference/id-prefixes#remove-prefixes) section of the Property Manager API (PAPI) documentation.

## Export a policy set and related policies

You can use the [Terraform CLI](https://github.com/akamai/cli-terraform) to export an existing policy set and its related policies into JSON files or directly into HCL syntax assuming that there are six or less levels of nested transformations.

You need to run the CLI separately for each policy set you want to add to your Terraform configuration.

 Running the CLI on a policy set also generates the resources for related policies and the JSON files of policies.
 If exported as schema, running the CLI on the policy set will not generate JSON files, but will generate the necessary data sources.

~> **Note:** If you use the [Image and Video Manager API](https://techdocs.akamai.com/ivm/reference/api), you can also modify JSON files you have for existing policy sets and policies.
In Control Center, you can view and download policy JSON files by clicking **View Policy JSON** in the Policy Editor.

## Add policy sets

Use the [`akamai_imaging_policy_set`](../resources/imaging_policy_set.md) resource to add a policy set to your Akamai Provider configuration.
Each policy set you're adding needs to have its own resource in your Terraform configuration.

You need the name of the policy set's JSON file that you created during the [export process](#export-policy-sets).

By default, when you save a policy set, it's automatically activated on both staging and production networks.
Nothing happens to the policy set until you link it to a property.

Also, when you create a policy set, it automatically creates a default `.auto` policy.
This policy provides the baseline settings that
determine how the policy generates derivative images and videos.
You can't delete the `.auto` policy without deleting the policy set.

## Add policies

Use these policy resources to add policies to your Akamai Provider configuration:

* **Image policies:** See the [imaging_policy_image](../resources/imaging_policy_image.md) resource.
* **Video policies:** See the [akamai_imaging_policy_video](../resources/imaging_policy_video.md) resource.

To set up these resources, you need the policy and policy set IDs found in the policy's JSON file.
This is one of the files created during the [export process](#export-policy-sets).
These IDs help link policy sets and policies to your property configuration.

Add a separate resource for each policy you're adding.

## Update the property rule tree

You need to add one or more imaging policy behaviors to the JSON rule tree file in each property you're adding image policies to.
See the [Set up property rules](../guides/get_started_property#set-up-property-rules) in the [Property Provisioning Module Guide](../guides/get_started_property.md) for additional information.

When adding these behaviors make sure that you use the correct policy set ID from the policy resource files.

If you wish to customize the behavior settings, see the Property Manager API (PAPI) behavior for the type of policy you're adding:

* **Image policies:** See the [`imageManager`](https://techdocs.akamai.com/property-mgr/reference/latest-image-manager) behavior.
* **Video policies:**  See the [`imageManagerVideo`](https://techdocs.akamai.com/property-mgr/reference/latest-image-manager-video) behavior.

### Example behaviors

**imageManager behavior:**

```
{
    "name":"imageManager",
    "options":{
        "enabled":true,
        "resize":false,
        "applyBestFileType":true,
        "cpCodeOriginal":{
            "id":6789
        },
        "cpCodeTransformed":{
            "id":12345
        },
        "policyTokenDefault":"default",
        "superCacheRegion":"US",
        "useExistingPolicySet":false,
        "advanced":false
    }
}
```

**imageManagerVideo behavior:**

```
{
    "name":"imageManagerVideo",
    "options":{
        "enabled":true,
        "resize":false,
        "applyBestFileType":true,
        "policyTokenDefault":"testToken",
        "superCacheRegion":"US",
        "useExistingPolicySet":false,
        "advanced":false
    }
}
```

## Activate the property on staging

To activate the property containing your newly added imaging policies, deploy the property to staging first.
In the [akamai_property_activation](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/property_activation) resource, set the `network` argument to `staging` and deploy the change.

~> **Note:** By default, Image and Video Manager policies and policy sets are automatically available on staging once you save them.

```
resource "akamai_property_activation" "example_staging" {
     property_id = akamai_property.example.id
     contact     = [local.email]
     # NOTE: Specifying a version as shown here will target the latest version created. This latest version will always be activated in staging.
     version     = akamai_property.example.latest_version
     network     = "STAGING"
     note        = "Sample activation"
}
```

## Test your images and videos

Now that the property, policy sets, and policies are active on the staging network, you should verify that your images and videos are appearing as expected.

See [Test your policy on staging](https://techdocs.akamai.com/ivm/docs/test-on-staging) for more information.

## Activate the policy on production

Once you finish testing on staging, activate the policy on production by setting `activate_on_production` to `true` in the policy resource.

If you're modifying an existing policy resource that's been activated in production, set `activate_on_production` to `false`.
Changing this value doesn't remove the policy in production, but lets you safely make changes to the policy in staging.

## Activate the property on production

Once the policy is on production, you can then deploy your updated property to production.
Update the [akamai_property_activation](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/property_activation) resource again, this time by setting the `network` argument to `production` and deploying the change.

## Maintenance

If needed, you can destroy policy sets and related policies.
When you remove a policy set, it removes all policies within that set including the `.auto` policy.

You can also remove individual policies from a policy set.
The exception is the `.auto` policy, which you can only remove if you delete the policy set.