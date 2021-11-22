---
layout: "akamai"
page_title: "Module: Application Security"
description: |-
   Application Security module for the Akamai Terraform Provider
---

# Application Security Module Guide

Application Security (appsec) in the Akamai Terraform provider (provider) enables application security configurations including such things as:

- Custom rules
- Match targets
- Other application security resources that operate within the cloud

This guide is for developers who:

- Are interested in implementing or updating an integration of Akamai functionality with Terraform.
- Already have some familiarity with Akamai products and Akamai application security.



------

### Before You Begin

This guide assumes that you have a basic understanding of Terraform and how it works (that is, you know how to install Terraform and the Akamai provider, how to configure your authentication credentials, how to create and use a .Terraform configuration file, etc.). If that’s not the case we strongly recommend you read the following two guides before going any further:

- [Akamai Provider: Get Started](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider)
- [Akamai Provider: Set Up Authentication](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/akamai_provider_auth)



------

## <a id="contents"></a>Table of Contents

- [Create a security configuration](#create)
- [Add a hostname to a security configuration](#hostname)
- [Activate a security configuration](#activate)
- [Create security policies by using Terraform](#policy)
- [Create a rate policy](#rate)
- [Create a match target](#match)
- [Modify a Kona rule set rule action](#kona)
- [Import a Terraform resource from one security configuration to another](#import)
- [Create an automated attack groups (AAG) security configuration](#aag)



------

## <a id="create"></a>Create a security configuration

[Back to table of contents](#contents)

[Create a blank security configuration](#blank)
[Create a security configuration that uses the recommended presets](#presets)
[Clone a security configuration](#clone)
[Activate a security configuration](#activateconfig)

In the Akamai world, a security configuration primarily functions as a container for other application security objects: security policies, rate policies, reputation profiles, Kona Site Defender rules, etc. Depending on your needs, you might have multiple security configurations; for example, many organizations have different security configurations for their different business units, their different domains, their different geographic units, and so on. If you need multiple configurations then you're in luck: as shown above, Terraform provides three different ways for you to create new security configurations.

### <a id="blank"></a>Create a blank security configuration

In Control Center, the easiest way to create a security configuration is to select the option **Manually (Create a blank Web Security Configuration**). This approach is easy because you don't have to do much beyond specifying:

- The contract and group ID associated with the new configuration.
- A name and (optionally) a description of the new configuration.
- At least one “selectable hostname.” (You can't create a security configuration without including at least one host to be protected by that configuration.)

This is also how you can do things in Terraform. If you go this route, just keep in mind that your configuration really *will* be empty: the new configuration will contain at least one host, but that's it. To make your configuration useful, you'll need to add security policies, match targets, rate policies, and so on. Admittedly, that can be a lot of work. However, If you want to start with the basics and then add on your new configuration from there, creating a blank (empty) security configuration doesn't require much effort on your part. In fact, the following Terraform configuration is all that's needed to create a blank configuration named **Empty Security Configuration**:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

resource "akamai_appsec_configuration" "create_config" {
 name        = "Empty Security Configuration"
 description = "This security configuration does not contain any items other than one hostname."
 contract_id = "1-3UW382"
 group_id    = 13139
 host_names  = ["documentation.akamai.com"]
}

output "create_config_id" {
 value = akamai_appsec_configuration.create_config.config_id
}
```

The first two blocks in this configuration are probably familiar to you: the initial block declares the Akamai Terraform provider, and the next block provides our authentication credentials. That sets the stage for the following:

```
resource "akamai_appsec_configuration" "create_config" {
 name        = "Empty Security Configuration"
 description = "This security configuration does not contain any items other than one hostname."
 contract_id = "1-3UW382"
 group_id    = 13139
 host_names  = ["documentation.akamai.com"]
}
```

This is where we use the akamai_appsec_configuration resource to create the new configuration. Creating a blank configuration requires us to supply the following five arguments:

| **Argument**  | **Description**                                              |
| ------------- | ------------------------------------------------------------ |
| `name`        | Unique name to be assigned to the new configuration.        |
| `description` | Brief description of the configuration and its intended purpose. |
| `contract_id` | Akamai contract ID associated with the new configuration. Use the [akamai_appsec_contracts_groups](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_contracts_groups) data source to return information about the contracts and  groups available to you. |
| `group_id`    | Akamai group ID associated with the new configuration.       |
| `host_names`  | Names of the selectable hosts to be protected by the configuration. Note that host names are passed as an array; that's what the square brackets surrounding  **"documentation.akamai.com"** are for. To add multiple hostnames to the configuration, separate the individual names by using commas. For example:     <br /><br />`host_names = ["documentation.akamai.com", "training.akamai.com", "events.akamai.com"]` |

Finally, we use this block to echo back the ID of the newly created security configuration:

```
output "create_config_id" {
 value = akamai_appsec_configuration.create_config.config_id
}
```

When we call the preceding configuration file we get a new security configuration as well as output similar to the following:

```
akamai_appsec_configuration.create_config: Creating...
akamai_appsec_configuration.create_config: Creation complete after 5s [id=76967]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

create_config_id = 71219
```



### <a id="presets"></a>Create a security configuration that uses the recommended presets

In some cases, you might want to get a security configuration up and running as quickly as possible, then fine-tune and adjust that configuration later. One easy way to do that is to create a security configuration that uses the recommended (default) settings. When you take this route, you end up with a security configuration that:

- Has a single security policy (a policy that employs the default values).
- Has the IP/Geo Firewall enabled. (Although that's only a convenience, because you won't have any network lists allowed or blocked by the firewall.)
- Has three rate policies: **Origin Server**, **POST Page Requests**, and **Page View Requests**. In all three policies, both the IPv4 and the IPv6 actions are set to **alert**.
- Has slow POST protection enabled using the default values. This includes setting the slow POST action to **alert**.
- Doesn't have any custom rules.
- Enables you to choose between automated attack groups (AAG, in which firewall rules are automatically updated and maintained for you), and Kona Rule Set (KRS) rules, which give you the ability to configure each rule's action, conditions, and exceptions.
- Has API request constraints enabled and the action set to **alert**. However, no API matches are defined.
- Has a single match target: **All Hostnames** (which matches the path **/***).

In this documentation, we won't create a configuration that includes all the default items. Instead, we'll create a simpler configuration that: 1) includes a security policy that uses the default security policy settings, and 2) sets the Web Application Firewall (WAF) mode to **KRS** . Although not a fully fleshed out configuration, this should give you an idea on how to:

- Create a security configuration.
- Use the ID of the new configuration to create a security policy.
- Use the ID of the security policy to modify the WAF mode.

Here's our sample configuration file:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

resource "akamai_appsec_configuration" "create_config" {
 name        = "Documentation Security Configuration"
 description = "This security configuration is used by the documentation team."
 contract_id = "1-3CV382"
 group_id    = 47346
 host_names  = ["ldap.host1.akamai.com.edgesuite-staging.net"]
}

resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id              = akamai_appsec_configuration.create_config.config_id
 default_settings       = true
 security_policy_name   = "Documentation Security Policy"
 security_policy_prefix = "doc1"
}

resource "akamai_appsec_waf_mode" "waf_mode" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 mode = "KRS"
}
```

After declaring the Akamai provider and providing our authentication credentials, we use the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_configuration) resource and this block to create the security configuration:

```
resource "akamai_appsec_configuration" "create_config" {
 name        = "Documentation Security Configuration"
 description = "This security configuration is used by the documentation team."
 contract_id = "1-3CV382"
 group_id    = 47346
 host_names  = ["ldap.host1.akamai.com.edgesuite-staging.net"]
}
```

If this looks familiar, that shouldn't come as a surprise: it's the same Terraform configuration used to create an empty security configuration. The only difference is that, when we created our empty configuration, we stopped right there. This time we'll keep going.

After the security configuration is complete, and after we know the ID of that configuration, we can create our security policy. But how are we supposed to know the ID of a brand-new security configuration?

As it turns out, determining the ID of a new configuration is surprisingly easy. To begin with, after a configuration has been created, the ID of that configuration is available in an attribute named `config_id`; we know that because that information is included in the [akamai_appsec_configuration resource documentation](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_configuration)n.

On top of that, the first line in our resource block looks like this:

```
resource "akamai_appsec_configuration" "create_config" {
```

As you know, akamai_appsec_configuration is the name of the resource that creates a configuration. Meanwhile, `create_config` is an attribute that represents the newly-created configuration. That means we have three pieces of information to work with:

| **Item**                      | **Description**                                              |
| ----------------------------- | ------------------------------------------------------------ |
| `akamai_appsec_configuration` | The Terraform resource.                                     |
| `create_config`               | The attribute that references the new configuration.        |
| `config_id`                   | The attribute that contains the ID of the new configuration. |

To reference the ID of our new configuration, we use “dot notation” to string these three items together:

```
akamai_appsec_configuration.create_config.config_id
```

If you look at the Terraform block that creates a security policy, you'll see that we use the preceding string as the value of the `config_id` argument:

```
resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id              = akamai_appsec_configuration.create_config.config_id
 default_settings       = true
 security_policy_name   = “"Documentation Security Policy"
 security_policy_prefix = "doc1"
}
```

> **Note**. We won't go into any of the other details involved in creating a security policy. You can find that information in the article **Creating a Security Policy**.

After we have a new configuration and a new security policy, all that's left is to specify the WAF mode. To do that, we need two things: the configuration ID and the security policy ID. We already have the configuration ID—we just used it to create the security policy. But what about the security policy ID?

As it turns out, any time you create a security policy an attribute named **security_policy_id** is made available to you. If you look at the first line of our security policy block, you'll see both a resource name (akamai_appsec_security_policy) and an attribute that represents the new policy (`security_policy_creat`e). How do we know the ID of our new security policy? Once again, we use dot notation to string together our individual elements:

```
akamai_appsec_security_policy.security_policy_create.security_policy_id
```

That's the value we assign to the `security_policy_id` argument:

```
resource "akamai_appsec_waf_mode" "waf_mode" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 mode               = "KRS"
}
```

All that's left now is to set the mode to KRS.

And that's it (at least for our immediate purposes). When we run the terraform apply command we'll get back a response similar to the following:

```
akamai_appsec_configuration.create_config: Creating...
akamai_appsec_configuration.create_config: Still creating... [10s elapsed]
akamai_appsec_configuration.create_config: Still creating... [20s elapsed]
akamai_appsec_configuration.create_config: Creation complete after 21s [id=76984]
akamai_appsec_security_policy.security_policy_create: Creating...
akamai_appsec_security_policy.security_policy_create: Creation complete after 7s [id=76984:doc1_137405]
akamai_appsec_waf_mode.waf_mode: Creating...
akamai_appsec_waf_mode.waf_mode: Creation complete after 4s [id=76984:doc1_137405]

Apply complete! Resources: 3 added, 0 changed, 0 destroyed.
```

If you take a closer look at the response, you'll see the ID of the new configuration (**76984**) and the ID of the new security policy (**doc1_137405**).



### <a id="clone"></a>Clone a security configuration

Another way to create a new—and functional—security configuration is to clone an existing configuration. In effect, you take security configuration A and create a (nearly) exact replica of that configuration (security configuration B). Among the items in configuration A replicated in configuration B are:

- Security policies
- Rate policies
- Custom rules
- Custom denies
- SIEM settings
- Advanced logging and prefetch settings
- Slow POST settings
- Match targets

In most cases items are copied exactly the way they appear in configuration A. For example, suppose configuration A contains the following security policies:

- Security policy 1
- Security policy 2
- Security policy 3
- Security policy 4

After the cloning operation is complete, configuration B contains the exact same set of policies, with the exact same names and setting values:

- Security policy 1
- Security policy 2
- Security policy 3
- Security policy 4

In other cases the Akamai provider takes into account the fact that the new configuration is a, well, new configuration. For example, in configuration A, the Web Security Configuration ID setting is shown as **90013**. When configuration A is cloned, however, the Web Security Configuration ID setting reflects the ID of the new configuration (configuration B).

A Terraform configuration for cloning a security configuration looks similar to this:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

resource "akamai_appsec_configuration" "clone_config" {
 name                  = "Cloned Security Configuration"
 description           = "This security configuration is based on configuration ID 90013."
 create_from_config_id = 90013
 create_from_version   = 9
 contract_id           = "1-3UW382"
 group_id              = 13139
 host_names            = ["documentation.akamai.com"]
}

output "clone_config_id" {
 value = akamai_appsec_configuration.clone_config.config_id
}
```

There's nothing especially complicated about this configuration. The configuration starts by declaring the Akamai Terraform provider, and by presenting our authentication credentials. From there, we move directly into the block that uses the akamai_appsec_configuration resource to clone configuration 90013:

```
resource "akamai_appsec_configuration" "clone_config" {
 name                  = "Cloned Security Configuration"
 description           = "This security configuration is based on configuration ID 90013."
 create_from_config_id = 90013
 create_from_version   = 9
 contract_id           = "1-3UW382"
 group_id              = 13139
 host_names            = ["documentation.akamai.com"]
}
```

The akamai_appsec_configuration block for cloning a configuration is remarkably similar to the other blocks we've looked at (such as the block for creating a blank configuration). The only difference is that when cloning a configuration we include two additional arguments:

```
create_from_config_id = 90013
create_from_version   = 9
```

The `create_from_config_id` argument specifies the ID of the security configuration you want to replicate (configuration A). Similarly, the `create_from_version` enables you to select a specific version of the configuration to be cloned. Note that this argument is optional. If it's not included, Terraform replicates the initial (v1) version of the configuration.

The final block of code is also optional. It simply echoes back the ID of the new security configuration:

```
output "clone_config_id" {
 value = akamai_appsec_configuration.clone_config.config_id
}
```

When you run the `terraform apply` command, you see output similar to this:

```
akamai_appsec_configuration.clone_config: Creating...
akamai_appsec_configuration.clone_config: Still creating... [10s elapsed]
akamai_appsec_configuration.clone_config: Creation complete after 11s [id=76982]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

clone_config_id = 76982
```

In this example, the new configuration has been assigned the ID **76982**.



### <a id="activateconfig"></a>Activate a security configuration

You can use Terraform to activate a configuration on either the staging network or the production network. A Terraform configuration for doing this looks similar to the following:

```
terraform {
 required_providers {
  akamai  = {
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

resource "akamai_appsec_activations" "activation" {
 config_id           = data.akamai_appsec_configuration.configuration.config_id
 network             = "STAGING"
 notes               = "This is a test configuration used by the documentation team."
 notification_emails = ["gstemp@akamai.com"]
}
```

For the most part, this is a typical Terraform configuration—we declare the Akamai provider, provide our authentication credentials, and connect to the Documentation configuration. We then use the [akamai_appsec_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_activations) resource and the following block to activate that configuration:

```
resource "akamai_appsec_activations" "activation" {
 config_id           = data.akamai_appsec_configuration.configuration.config_id
 network             = "STAGING"
 notes               = "This is a test configuration used by the documentation team."
 notification_emails = ["gstemp@akamai.com"]
}
```

Inside this block we include these arguments and argument values:

| **Argument**          | **Description**                                              |
| --------------------- | ------------------------------------------------------------ |
| `config_id`           | Unique identifier of the configuration being activated.     |
| `network`             | Name of the network the configuration is being activated for. Allowed values are:     <br />*  **staging**  <br />*  **production** |
| `notes`               | Information about the configuration and the reason for its activation. |
| `notification_emails` | JSON array of email addresses of people to be notified when activation is complete. To send notification emails to multiple people, separate the individual email addresses by using commas:     <br /><br />`notification_emails = ["gstemp@akamai.com", "karim.nafir@mail.com"]` |

From here we can run `terraform plan` to verify our syntax, then run `terraform apply` to activate the security configuration.



------

## <a id="hostname"></a>Add a hostname to a security configuration

[Back to table of contents](#contents)

[Return a collection of selectable hostnames](#return)
[Add a hostname to a security configuration](#addhost)

Security configurations are not designed to automatically protect every server in your organization. Instead, each security configuration protects only those servers (that is, only those hosts) designated for this protection. Do you have different servers that require different levels or different types of protection? That's fine: put one set of servers in Configuration A, another set in Configuration B, and so on. Doing this gives you the ability to fine-tune your websites and to strike a balance between security and efficiency.

So how do you designate a host for protection by a security configuration? That’s a two-step process:

1. **Determine which hosts are selectable—that is, which hosts are available to be added to a security configuration**. Note that it's possible that some of your hosts *can't* be added to a security configuration—for example, if a contract has expired you won’t be able to offer protection to certain servers. Note as well that a single security configuration can manage multiple hosts. However, an individual host can be a member of only one security configuration. In other words, Host A can belong to Configuration A or Configuration B, but it can't belong to both Configuration A and Configuration B at the same time.
2. **Add the hostname to the appropriate security configuration**. When updating the selected hosts for a security configuration you can either append the new host (or hosts—you can add multiple hosts in a single operation) to the existing collection, or you can replace the existing collection with the hosts specified in your Terraform configuration.

When multiple adding hosts you can specify each host individually or you can use wildcard characters. For example, this syntax adds four hosts (host1, host2, host3, and host4) from the akamai.com domain to a security configuration:

```
hostnames = ["host1.akamai.com", "host2.akamai.com", "host3.akamai.com", "host4.akamai.com"]
```

That syntax works fine. But so does this syntax, which uses a wildcard character (*****) to add all of your akamai.com hosts to the configuration:

```
hostnames = [ "*.akamai.com"]
```



### <a id="return"></a>Return a collection of selectable hostnames

To return a list of available hostnames, use the [akamai_appsec_selectable_hostnames](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_selectable_hostnames) data source. There are two ways to use this data source: you can return all the available hostnames for a specific contract and group, or you can return all the available hostnames for a specific configuration. For example, this Terraform configuration returns all the selectable hostnames for contract **7-5DR99** and group **57112**:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames_for_create_configuration" {
 contractid = "7-5DR99"
 groupid    = 57112
}

output "selectable_hostnames_for_create_configuration" {
 value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.hostnames
}
```

Alternatively, this configuration returns the selectable hostnames for configuration **90013**:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
 config_id = 90013
}

output "selectable_hostnames" {
 value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}
```

Either way, if you call this configuration by using the `terraform plan` command, you'll get back a list of hostnames similar to the following:

```
Changes to Outputs:

selectable_hostnames_for_create_configuration = [
 "host1.akamai.com",
 "host2.akamai.com",
 "host3.akamai.com",
 "host4.akamai.com"
]
```

Any of these hostnames can be added to your configuration.

Incidentally, if you try adding a hostname that *isn't* on the list of selectable hostnames then your command fails with an error similar to this:

```
│ Error: setting property value: Title: Invalid Input Error; Type: https://problems.luna.akamaiapis.net/appsec-configuration/error-types/INVALID-INPUT-ERROR; Detail: You don't have access to add these hostnames [identitydocs.akamai.com]
│
│  with akamai_appsec_selected_hostnames.appsecselectedhostnames,
│  on akamai.tf line 15, in resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames":
│  15: resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
```



### <a id="addhost">Add a hostname to a security configuration

After you've chosen a hostname (or set of hostnames), you can add these hosts to a security configuration by using the [akamai_appsec_selected_hostnames](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_selected_hostnames) resource and a Terraform configuration like this one:

```
terraform {
 required_providers {
  akamai  = {
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

resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
 config_id = data.akamai_appsec_configuration.configuration.config_id
 hostnames = [ "host1.akamai.com", "host2.akamai.com"]
 mode      = "APPEND"
}
```

In this configuration, we start by:

1. Declaring the Akamai provider.
2. Providing our authentication credentials.
3. Connecting to the **Documentation** security configuration.

When all that's done, we use this block of code to add the hostnames **host1.akamai.com** and **host2.akamai.com** to the **Documentation** configuration:

```
resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
 config_id = data.akamai_appsec_configuration.configuration.config_id
 hostnames = [ "host1.akamai.com", "host2.akamai.com"]
 mode      = "APPEND"
}
```

In order to add the hostnames to the existing collection of hostnames for the configuration, we set the `mode` to **APPEND**. For example, suppose our collections currently has these three hosts:

- hostA.akamai.com
- hostB.akamai.com
- hostC.akamai.com

After we run our Terraform configuration the security configuration has these hosts:

- hostA.akamai.com
- hostB.akamai.com
- hostC.akamai.com
- host1.akamai.com
- host2.akamai.com

By comparison, if we set `mode` to **REPLACE**, the existing set of hostnames are deleted and are *replaced* by the hostnames specified in the Terraform configuration. In other words, the security configuration ends up with these hosts:

- host-1.akamai.com
- host-2.akamai.com





------

## <a id="activate"></a>Activate a security configuration

[Back to table of contents](#contents)

[Reactivate a security configuration](#reactivate)

Security configurations must be activated before they can begin analyzing and responding to user requests. Typically, activation is a two-step process: first the configuration is activated on the staging network and then, after testing and fine-tuning, the configuration is activated on the production network. At that point, the configuration is fully deployed, and *is* analyzing and responding to requests.

When the time comes, you can use Terraform to activate a configuration on either the staging network or the production network. A Terraform configuration for doing this looks similar to the following:

```
terraform {
 required_providers {
  akamai  = {
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

resource "akamai_appsec_activations" "activation" {
 config_id           = data.akamai_appsec_configuration.configuration.config_id
 network             = "STAGING"
 activate            = true
 notes               = "This is a test configuration used by the documentation team."
 notification_emails = ["gstemp@akamai.com"]
}
```

For the most part, this is a typical Terraform configuration—we declare the Akamai provider, provide our authentication credentials, and connect to the **Documentation** configuration. We then use the [akamai_appsec_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_activations) resource and the following block to activate that configuration:

```
resource "akamai_appsec_activations" "activation" {
 config_id           = data.akamai_appsec_configuration.configuration.config_id
 network             = "STAGING"
 activate            = true
 notes               = "This is a test configuration used by the documentation team."
 notification_emails = ["gstemp@akamai.com"]
}
```

Inside this block we include the following arguments and argument values:

| **Argument**          | **Description**                                              |
| --------------------- | ------------------------------------------------------------ |
| `config_id`           | Unique identifier of the configuration being activated.     |
| `network`             | Name of the network the configuration is being activated on. Allowed values are:     <br />*  **staging**  <br />*  **production** |
| `activate`            | If **true** (the default value), the security configuration will be activated; if  **false**, the security configuration will be *deactivated*. Note that this argument is optional—if it's not included the security configuration will be activated. |
| `notes`               | Information about the configuration and its activation.     |
| `notification_emails` | JSON array of email addresses of the people to be notified when activation is complete. To send notification emails to multiple people, separate the individual email addresses by using commas:<br /><br />`notification_emails = ["gstemp@akamai.com", "karim.nafir@mail.com"]` |

From here we can run `terraform plan` to verify our syntax, then run `terraform apply` to activate the security configuration. If everything goes as expected, you'll see output similar to the following:

```
akamai_appsec_activations.activation: Creating...
akamai_appsec_activations.activation: Creation complete after 2s [id=none]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```



### <a id="reactivate"></a>Reactivate a security configuration

Depending on the changes you make to it, a security configuration might need to be reactivated at some point. However, if you use the exact same Terraform syntax previously employed in the previous example, reactivation *doesn’t* take place. Instead, Terraform doesn’t do anything at all:

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

As you can see, no resources are added, changed, or destroyed. In other words, nothing happens.

The problem here lies in the way that Terraform processes .tf files. When we originally activated the security configuration, we used the Terraform block shown in the previous section of the documentation. When we try to reactivate the configuration using that same code, Terraform is unable to see any changes: the `config_id` is the same, the `network` is the same, etc. Because nothing seems to have changed, Terraform does just that: nothing.

Perhaps the best way to work around this issue is to change the value of the notes argument: this enables you to make a change of *some* kind without having to make a more drastic change (for example, changing the network or the notification list). When we originally activated the configuration we used this line:

```
notes = "This is a test configuration used by the documentation team."
```

To reactivate the security configuration, we simply need to change the value to something (anything) else:

```
notes = "This is a reactivated test configuration."
```

Now you can reactivate the configuration. If you need to run activation again, just make another change to the notes argument.



------

## <a id="policy"></a>Create security policies by using Terraform

[Back to table of contents](#contents)

[A note about using multiple security policies](#multiple)
[Security policy default settings](#defaultsettings)
[Create a security policy that uses the default settings](#createdefault)
[Clone a security policy](#clonepolicy)
[Add custom setting values when you create a security policy](#customvalues)

At heart, websites are all about handling user requests: this user would like to download a file; that user would like to stream a video; a third user would like to visit one of your web pages. The vast majority of these requests are legitimate, and harmless; in fact, they're the very reason you published your website in the first place. However, other requests (either maliciously or inadvertently) might not be so harmless. To safely and securely manage your website, you need to be able to identify suspect requests and quickly and efficiently deal with those requests.

At Akamai, security policies play a key role in identifying and handling website requests. If a request is flagged by a match target (that is, if the request matches criteria you have specified in advance) the security policy associated with that match target can step in and provide a more detailed analysis on the request, applying protections such as rate limiting and reputation controls to help verify the legitimacy and the safety of the request. Requests that pass these tests are allowed through; depending on how you have configured your policies, requests that don't pass these tests can be rejected.

That, in a nutshell, is why you need security policies.


### <a id="multiple"></a>A note about using multiple security policies

A single security configuration can have multiple security policies. And that's good: after all, you might, to list one example, have one set of APIs that require a different set of protections than your other APIs. In that case, you might need two security policies: one for the “special” set of APIs and the other for the remaining APIs.

What you probably *don't* need, however, is one security policy for each individual API. Having multiple security policies provides you with flexibility and with the opportunity to fine-tune your protections. At the same time, however, each security policy you add to a security configuration increases the time it takes to analyze and process each request. You'll need to find a balance between having customized protections and having an efficient and responsive website. As a general rule, the fewer security policies you employ the better.



### <a id="defaultsettings"></a>Security policy default settings

Unless you specify otherwise, any new security policy you create is assigned the default policy settings. The current values for these settings can be returned by using the [akamai_appsec_security_policy_protections](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_security_policy_protections) data source and a Terraform configuration similar to this:

```
terraform {
 required_providers {
  akamai  = {
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

data "akamai_appsec_security_policy_protections" "protections" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637"
}

output "protections_response" {
 value = data.akamai_appsec_security_policy_protections.protections.output_test
}
```

As you can see, this is a pretty standard configuration file: we call the Akamai provider and point that provider to our authentication credentials (stored in the .edgerc file), We connect to the **Documentation** configuration, then use this block to return protections information for the security policy **gms1_134637**:

```
data "akamai_appsec_security_policy_protections" "protections" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637"
}
```

After the protections data is returned, we use the final block in the file to output that information to the screen:

```
output "protections_response" {
 value = data.akamai_appsec_security_policy_protections.protections.output_test
}
```

When all is said and done you'll get back response similar to this:

```
+------------------------------------------------------------------------------------------------------------------------------------------+
| wafProtectionDS                                                             |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| APICONSTRAINTS | APPLICATIONLAYERCONTROLS | BOTMANCONTROLS | NETWORKLAYERCONTROLS | RATECONTROLS | REPUTATIONCONTROLS | SLOWPOSTCONTROLS |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| true      | true           | false     | true         | true     | true        | true       |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
```

These properties are described in the following table:

| **API property**           | **Description**                                              | **Default value** |
| -------------------------- | ------------------------------------------------------------ | ----------------- |
| `APICONSTRAINTS`           | Places limits on both the number and the size of API requests sent by a given user. | true              |
| `APPLICATIONLAYERCONTROLS` | Uses the Web Application Firewall (WAF) to help minimize the effects of cross-site  scripting, SQL injection, file inclusion, and other attacks. | true              |
| `BOTMANCONTROLS`           | Places limits on the number of valid, and invalid, form submissions associated with a single user. | false             |
| `NETWORKLAYERCONTROLS`     | Blocks (or allows) requests based on a client's IP address or geographic location. | true              |
| `RATECONTROLS`             | Provides  a way to monitor, and to control, the rate of requests received by your site. | true              |
| `REPUTATIONCONTROLS`       | Helps identify potentially-malicious clients based on past behaviors associated with the client IP address. | false             |
| `SLOWPOSTCONTROLS`         | Helps guard against Denial of Service attacks caused by extremely slow request rates. | true              |

When you enable one of these features, that feature is automatically assigned a set of default values. For example, if you enable slow POST controls then, by default, the feature is configured with the following settings:

| **Action** | **SLOW_RATE_THRESHOLD** **RATE** | **SLOW_RATE_THRESHOLD** **PERIOD** | **DURATION_THRESHOLD** **TIMEOUT** |
| ---------- | --------------------------------- | ----------------------------------- | ----------------------------------- |
| alert      | 10                                | 60                                  | Null                                |

If those settings work for you that's great. If they don't, you can use an additional terraform block to change the setting values as needed. For example, this Terraform snippet modifies all four of the slow POST property values:

```
data "akamai_appsec_configuration" "configuration" {
 name = var.security_configuration
}

resource "akamai_appsec_slow_post" "slow_post" {
 config_id                  = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id         = var.security_policy_id
 slow_rate_action           = "abort"
 slow_rate_threshold_rate   = 15
 slow_rate_threshold_period = 30
 duration_threshold_timeout = 20
}
```



### <a id="createdefault"></a>Create a security policy that uses the default settings

The quickest and easiest way to create a security policy is to create a policy that uses the default settings. As a reminder, those defaults are shown below:

| **Setting**                | **Default Value** |
| -------------------------- | ----------------- |
| `APICONSTRAINTS`           | true              |
| `APPLICATIONLAYERCONTROLS` | true              |
| `BOTMANCONTROLS`           | false             |
| `NETWORKLAYERCONTROLS`     | true              |
| `RATECONTROLS`             | true              |
| `REPUTATIONCONTROLS`       | false             |
| `SLOWPOSTCONTROLS`         | true              |

Here's a sample Terraform collection that creates a security policy that uses those default settings:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
 }

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
 name = "Documentation”
}

resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id              = data.akamai_appsec_configuration.configuration.config_id
 default_settings       = true
 security_policy_name   = "New Default Policy"
 security_policy_prefix = "gms2"
}

output "security_policy_create" {
 value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

As you can see, this is another straightforward Terraform configuration: it simply creates a new security policy (with the name **New Default Policy** and the prefix **gms2**) and associates that policy with the **Documentation** security configuration. Note that you don't include an ID for the security policy when creating your configuration file: IDs are assigned by the system when the policy is created. That ID comprises the security policy prefix, an underscore (**_**) and a numeric value assigned by Akamai. For example:

```
gms2_135566
```

 Note, too that our Terraform configuration also includes this argument:

```
default_settings = true
```

This tells Terraform that we want the new security policy to use the default security policy settings. The `default_settings` argument optional: if it's not included the policy will automatically be assigned the default settings. And if you set `default_settings` to **false**? In that case, all the policy settings are also set to false:

```
+------------------------------------------------------------------------------------------------------------------------------------------+
| wafProtectionDS                                                             |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| APICONSTRAINTS | APPLICATIONLAYERCONTROLS | BOTMANCONTROLS | NETWORKLAYERCONTROLS | RATECONTROLS | REPUTATIONCONTROLS | SLOWPOSTCONTROLS |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| false     | false          | false     | false        | false    | false        | false     |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
```

 After your Terraform configuration is ready, run the `terraform plan` command from the command prompt. Running the `plan` command does two things for you. First, it does some syntax checking, and alerts you to many (although not necessarily all) configuration errors:

```
│ Error: Missing required argument
│
│  on akamai.tf line 17, in resource "akamai_appsec_security_policy" "security_policy_rename":
│  17: resource "akamai_appsec_security_policy" "security_policy_rename" {
│
│ The argument "security_policy_prefix" is required, but no definition was found.
```

 Second, it tells you exactly what happens if you create the security policy:

```
Terraform will perform the following actions:

akamai_appsec_security_policy.security_policy_create will be created

resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id              = 90013
 default_settings       = true
 id                     = (known after apply)
 security_policy_id     = (known after apply)
 security_policy_name   = "Default Settings Policy"
 security_policy_prefix = "gms2"
}

Plan: 1 to add, 0 to change, 0 to destroy.
Changes to Outputs:
security_policy_create = (known after apply)
```

 If everything looks OK then run `terraform apply` from the command prompt (and answer **yes** to the prompt that asks if you really want to apply these actions). If the policy can be created, you'll see output similar to this:

```
akamai_appsec_security_policy.security_policy_create: Creating...
akamai_appsec_security_policy.security_policy_create: Creation complete after 9s [id=90013:gms2_135566]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

security_policy_create = "gms2_135566"
```

 Incidentally, the very last line in the response (the line that echoes back the security policy ID) shows up because our configuration included this block:

```
output "security_policy_create" {
 value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

 If you leave out that block you'll still get back the security policy ID; it just won't be as easy to find:

```
akamai_appsec_security_policy.security_policy_create: Creating...
akamai_appsec_security_policy.security_policy_create: Creation complete after 9s [id=90013:gms2_135566]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

> **Hint**. It's the **gms2_135566** portion of **id=90013:gms2_135566**. That value happens to be the ID of the security configuration (**90013**) followed by a colon followed by the security policy ID.



### <a id="clonepolicy"></a>Clone a security policy

To clone something means to make an exact replica of that something. That's exactly what happens when you clone a security policy: you create a new policy (policy B) that has the exact same settings and setting values as an existing policy (policy A). Other than their identifiers (such as the policy name and the policy ID), the two policies will be indistinguishable from one another.

To explain how to clone a security policy, let's start by looking at a sample configuration file. As you'll see, the process used to clone a security policy is very similar to the process used to create a security policy that uses the default settings:

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
 name = "Documentation”
}

resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id                      = data.akamai_appsec_configuration.configuration.config_id
 default_settings               = false
 create_from_security_policy_id = "gms1_134637"
 security_policy_name           = "Cloned Policy"
 security_policy_prefix         = "gms3"
}

output "security_policy_create" {
 value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

The preceding configuration creates a new security policy (a policy with the name **Cloned Policy** and the prefix **gms3**) as part of the **Documentation** security configuration. Unlike our previous example, this new policy isn't based on the default settings; instead, the policy is, for all intents and purposes, a duplicate of the existing security policy **gms1_134637**. That means that the policy is configured using the exact same settings and setting values assigned to **gms1_134637**. In other words:

| **API property**           | **gms_134637** | **Cloned policy** |
| -------------------------- | -------------- | ----------------- |
| `APICONSTRAINTS`           | false          | false             |
| `APPLICATIONLAYERCONTROLS` | false          | false             |
| `BOTMANCONTROLS`           | false          | false             |
| `NETWORKLAYERCONTROLS`     | true           | true              |
| `RATECONTROLS`             | true           | true              |
| `REPUTATIONCONTROLS`       | false          | false             |
| `SLOWPOSTCONTROLS`         | true           | true              |

In order to clone a security policy, your Terraform configuration must include the following two arguments:

```
default_settings = false
create_from_security_policy_id = "gms1_134637"
```

The first argument (`default_settings = false`) tells Terraform *not* to apply the default settings to the new policy. Instead, we want the new policy to have the same settings and settings values as the ones assigned to the existing security policy **gms1_134637**. As you might have guessed, that's what the second argument does: it specifies the ID of the security policy whose settings we want copied to the new policy.

The rest is easy. Just like we did before, run `terraform plan` to verify your syntax, and then run `terraform apply` to create the new security policy. If everything goes according to plan, the policy is created, and you'll see output similar to this:

```
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
Outputs

security_policy_create = "gms3_135568"
```



### <a id="customvalues"></a>Add custom setting values when you create a security policy

As alluded to previously, another option available when creating a security policy is to create a policy that includes custom setting values; for example, you can create a policy that enables slow POST protection, but that doesn't use the default setting values:

| **Slow POST setting**       | **Default value** | **New policy value** |
| --------------------------- | ----------------- | -------------------- |
| action                      | alert             | abort                |
| SLOW_RATE_THRESHOLD  RATE   | 10                | 15                   |
| SLOW_RATE_THRESHOLD  PERIOD | 60                | 30                   |
| DURATION_THRESHOLD  TIMEOUT | null              | 20                   |

To do this, we'll use a Terraform configuration similar to this:

```
terraform {
 required _providers {
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

resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id              = data.akamai_appsec_configuration.configuration.config_id
 default_settings       = true
 security_policy_name   = "New Modified Policy"
 security_policy_prefix = "gms4"
}

output "security_policy_create" {
 value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}

resource "akamai_appsec_slowpost_protection" "protection" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}

resource "akamai_appsec_slow_post" "slow_post" {
 config_id                  = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id         = akamai_appsec_security_policy.security_policy_create.security_policy_id
 slow_rate_action           = "abort"
 slow_rate_threshold_rate   = 15
 slow_rate_threshold_period = 30
 duration_threshold_timeout = 20
}
```

The first half of this configuration should look familiar; it simply creates a new security policy (**New Modified Policy**) that uses the default settings:

```
default_settings       = true
security_policy_name   = "New Modified Policy"
security_policy_prefix = "gms4"
```

In this case, however, the configuration doesn't end after the policy has been created. Instead, it continues on to assign custom values to the slow POST control settings:

```
resource "akamai_appsec_slow_post" "slow_post" {
 config_id = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id         = akamai_appsec_security_policy.security_policy_create.security_policy_id
 slow_rate_action           = "abort"
 slow_rate_threshold_rate   = 15
 slow_rate_threshold_period = 30
 duration_threshold_timeout = 20
}
```

> **Note**. There was no need to enable slow POST controls because the policy was created using the default settings. By default, slow POST control is enabled in all new security policies.

In the preceding block, the only “tricky” line is this one:

```
security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
```

In this line we specify the ID of the security policy we want to update. Of course, this is a brand-new security policy which didn't even exist a few seconds ago. So how do we know the ID of the new policy? By using this value:

```
akamai_appsec_security_policy.security_policy_create.security_policy_id
```

This is the value we echo back to the screen immediately after the security policy has been created:

```
output "security_policy_create" {
 value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

That's also the ID of the new policy. Note the final line in this snippet from the Terraform response when creating a new policy:

```
akamai_appsec_security_policy.security_policy_create: Creating...
akamai_appsec_security_policy.security_policy_create: Creation complete after 8s [id=90013:gms4_135620]
```

In other words, what we do here is:

1. Create a new security policy, one that uses the default settings.
2. Retrieve the ID from that new policy and then use that ID to modify the slow POST control settings.

That's all there is to it. If everything works, we should see the following output as our Terraform configuration completes:

```
akamai_appsec_security_policy.security_policy_create: Creating...
akamai_appsec_security_policy.security_policy_create: Creation complete after 8s [id=90013:gms4_135620]
akamai_appsec_slow_post.slow_post: Creating...
akamai_appsec_slowpost_protection.protection: Creating...
akamai_appsec_slowpost_protection.protection: Creation complete after 4s [id=90013:gms4_135620]
akamai_appsec_slow_post.slow_post: Creation complete after 4s [id=90013:gms4_135620]

Apply complete! Resources: 3 added, 0 changed, 0 destroyed.

Outputs:

security_policy_create = "gms4_135620"
```




------

## <a id="rate"></a>Create a rate policy

[Back to table of contents](#contents)

[Configure rate policy actions](#actions)
[The rate policy JSON file](#ratejson)

A “classic” way to take down a website is to overwhelm the site with requests, transmitting so many requests that the site exhausts itself trying to keep up. This might be done maliciously (with the intent of crashing your web servers) or it might be done inadvertently: for example, if you announce a special offer to anyone who visits your site in the next hour you might get so many visitors that this legitimate traffic ends up bringing the site down.

In other words, and in some cases. it's possible to have too much of a good thing. Because of that, it's important that you monitor and moderate the number and rate of all the requests you receive. In the Akamai world, managing request rates is primarily done by employing a set of rate control policies. These policies revolve around two measures:

- **averageThreshold**. Measures the average number of requests recorded during a two-minute interval. The threshold value is the total number of requests divided by 2 minutes (120 seconds).
- **burstThreshold**. Measures the average number of requests recorded during a 5-second interval.

When configuring a rate policy, the average threshold should always be less than the burst threshold. Why? Well, the burst threshold often measures a brief flurry of activity that disappears as quickly as it appears. By contrast, the average threshold measures a much longer period of sustained activity. A sustained rate of activity obviously has the capability to create more problems than a brief and transient rate of activity.

Rate policies are also designed to trigger only when certain conditions are met. For example, a policy might be configured to fire only when a request results in a specified HTTP response code (e.g., a 404 or 500 error).

Terraform provides a way to quickly (and easily) create rate policies: this is done by specifying rate policy properties and property values in a JSON file and then running a Terraform configuration that creates a new policy based on those values. After a policy is created, you can use an additional Terraform block to assign an action to the policy (i.e., issue an alert if a policy threshold has been breached; deny the request if a policy threshold has been breached; etc.).

For example, the following Terraform configuration creates a new rate policy and assigns the policy to the **Documentation** security configuration. After the policy has been created, the configuration assigns a pair of actions to the new policy:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
 name = “Documentation”
}

resource "akamai_appsec_rate_policy" "rate_policy" {
 config_id   = data.akamai_appsec_configuration.configuration.config_id
 rate_policy = file("${path.module}/rate_policy.json")
}

output "rate_policy_id" {
 value = akamai_appsec_rate_policy.rate_policy.rate_policy_id
}

resource "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637"
 rate_policy_id     = akamai_appsec_rate_policy.rate_policy.rate_policy_id
 ipv4_action        = "alert"
 ipv6_action        = "alert"
}
```

As you can see, this Terraform configuration is similar to our other Akamai Terraform configurations. After identifying the provider and providing our credentials, we connect to the security configuration:

```
data "akamai_appsec_configuration" "configuration" {
 name = “Documentation”
}
```

We then use this block and the [akamai_appsec_rate_policy](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_policy) resource to create the policy:

```
resource "akamai_appsec_rate_policy" "rate_policy" {
 config_id   = data.akamai_appsec_configuration.configuration.config_id
 rate_policy = file("${path.module}/rate_policy.json")
}
```

If you're looking for all the configuration values for the new policy, you won't find them here; instead, those values are defined in a JSON file named **rate_policy.json**. (Incidentally, that file name is arbitrary: you can give your file any name you want). This argument tells Terraform that the properties and property values for the new rate policy should be read from rate_policy.json:

```
rate_policy = file("${path.module}/rate_policy.json")=
```

>  **Note**. We'll look at a sample JSON file in a few minutes.

The syntax **\${path.module}** is simply a shorthand way to specify that the JSON file is stored in the same folder as the Terraform executable. You don't have to store your JSON files in the same folder as the Terraform executable: just remember that, if you use a different folder, you'll need to specify the full path to the JSON file. Otherwise, Terraform won't be able to find it.

All that's left now is to echo back the ID of the new policy:

```
output "rate_policy_id" {
 value = akamai_appsec_rate_policy.rate_policy.rate_policy_id
}
```

Well, that and configure policy actions for the new policy.



### <a id="actions"></a>Configure rate policy actions

After you've created your rate policy, you'll want to configure the actions to be taken any time the policy is triggered (e.g., any time the `burstThreshold` is exceeded). There are four options available to you, and these options must be set for both IPv4 and IPv6 IP addresses:

- **alert**. An alert is issued if the policy is triggered.
- **deny**. The request is denied if the policy is triggered. It's recommended that you don't start out by setting a rate policy action to deny. Instead, start by setting all your actions to **alert** and then spend a few days monitoring and fine-tuning your policy threshold before you begin denying requests. If you don' do this, you run the risk of denying more requests than you really need to.
- **deny_custom_{custom_deny_id}**. Takes the action specified by the custom deny.
- **none**. No action of any kind is taken if the policy is triggered.

>  **Note**. As you'll see later in this documentation, rate policies have a property named **sameActionOnIpv** that indicates whether the same action (for example, **deny**) is used on both IPv4 and IPv6 addresses. When setting a rate policy action, however, you must specify both the IPv4 and IPv6 actions. For example, if you don't include the IPv4 action, then your configuration will fail because a required argument (in this example, `ipv4_action`) is missing.

In our sample configuration, we use the following Terraform block to set the rate policy actions:

```
resource "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637"
 rate_policy_id     = akamai_appsec_rate_policy.rate_policy.rate_policy_id
 ipv4_action        = "alert"
 ipv6_action        = "alert"
}
```

To set the actions, use the [akamai_appsec_rate_policy_action](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_policy_action) resource, and specify the appropriate security configuration (`config_id`), security policy (`security_policy_id`), and our newly-created rate policy (`rate_policy_id`). To indicate the rate policy, we reference the ID of that policy:

```
rate_policy_id = akamai_appsec_rate_policy.rate_policy.rate_policy_id
```

At that point all that's left is to configure the IPv4 and IPv6 actions:

```
ipv4_action = "alert"
ipv6_action = "alert"
```

Keep in mind that these two actions don't have to be the same. Although it's recommended that new rate policies have both these actions set to **alert**, you can specify that IPv4 addresses trigger a different action than IPv6 addresses:

```
ipv4_action = "deny"
ipv6_action = "alert"
```



### <a id="ratejson"></a>The rate policy JSON file

A JSON file used to define rate policy properties and property values looks similar to this:

```
{
  "additionalMatchOptions": [{
    "positiveMatch": true,
    "type": "ResponseStatusCondition",
    "values": ["400", "401", "402", "403", "404", "405", "406", "407", "408", "409", "410", "500", "501", "502", "503", "504"]
  }],
  "averageThreshold": 5,
  "burstThreshold": 8,
  "clientIdentifier": "ip",
  "description": "An excessive error rate from the origin could indicate malicious activity by a bot scanning the site or a publishing error. In both cases, this would increase the origin traffic and could potentially destabilize it.",
  "matchType": "path",
  "name": "HTTP Response Codes",
  "pathMatchType": "Custom",
  "pathUriPositiveMatch": true,
  "requestType": "ForwardResponse",
  "sameActionOnIpv6": true,
  "type": "WAF",
  "useXForwardForHeaders": false
}
```

Although the rate policies you create will use a JSON file *similar* to the one shown above, there will be differences depending on such things as your `matchType`, your `, etc. Rate policy properties available to you are briefly discussed in the following sections of the documentation.

#### Required properties

Any rate policy JSON file you create must include the properties shown in the following table:

| **Property**       | **Datatype** | **Description**                                              |
| ------------------ | ------------ | ------------------------------------------------------------ |
| `averageThreshold` | integer      | Maximum number of allowed hits per second during any two-minute interval. |
| `burstThreshold`   | integer      | Maximum number of allowed hits per second during any five-second interval. |
| `clientIdentifier` | string       | Identifier used to identify and track request senders; this value is required only when using Web Application Firewall. Allowed values are:     <br />*  **api-key**. Supported only for API match criteria.  <br />*  i**p-useragent**. Typically preferred over ip when identifying a client.     <br />*  **ip**. Identifies clients by IP address.  <br />*  **cookie:value**. Helps track requests over an individual session, even if the IP address changes. |
| `matchType`        | string       | Indicates the type of path matched by the policy allowed values are:     <br />*  **path**. Matches website paths.  <br />*  **api**. Matches API paths. |
| `name`             | string       | Unique name assigned to a rate policy.                      |
| `pathMatchType`    | string       | Type of path to match in incoming requests. Allowed values are:     <br />*  **AllRequests**. Matches an empty path or any path that ends in a trailing slash (**/**).  <br />*  **TopLevel**. Matches top-level hostnames only.  <br />*  **Custom**. Matches a specific path or path component. This property is only required when the `matchType` is set to **path**. |
| `requestType`      | string       | Type of request to count towards the rate policy's thresholds. Allowed values are:    <br />*  **ClientRequest**. Counts client requests to edge servers.  <br />*  **ClientResponse**. Counts edge responses to the client.  <br />*  **ForwardResponse**. Counts origin responses to the client.  <br />*  **ForwardRequest**. Counts edge requests to your origin. |
| `sameActionOnIpv6` | boolean      | Indicates whether the same rate policy action applies to both IPv6 traffic and IPv4 traffic. |
| `type`             | string       | Rate policy type. Allowed values are:     <br />*  **WAF**. Web Application Firewall.  <br />*  **BOTMAN**. Bot Manager. |

#### Optional properties

Optional rate policy properties are described in the following table:

| **Property**            | **Datatype** | **Description**                                              |
| ----------------------- | ------------ | ------------------------------------------------------------ |
| `description`           | string       | Descriptive text about the policy.                          |
| `hostnames`             | array        | Array of hostnames that trigger a policy match. If a hostname is not in the array  then that request is ignored by the policy. |
| pathUriPositiveMatch    | boolean      | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). |
| `useXForwardForHeaders` | boolean      | Indicates whether the policy checks the contents of the **X-Forwarded-For** header in incoming requests. |

**The additionalMatchOptions object**

Specifies additional matching conditions for the rate policy. For example:

```
"additionalMatchOptions": [
 {
  "positiveMatch": false,
  "values": [
   "121989_DOCUMENTATION001",
   "060389_DOCUMENTATION002"
   ],
  "type": "NetworkListCondition"
 }
]
```

Properties of the `addtionalmatchOptions` object are described in the following table:

| **Property**    | **Datatype** | **Description**                                              |
| --------------- | ------------ | ------------------------------------------------------------ |
| `properties`    | string       | Match condition type. Allowed values are:     <br />*  **IpAddressCondition**  <br />*  **NetworkListCondition**  <br />*  **RequestHeaderCondition**  <br />*  **RequestMethodCondition**  <br />*  **ResponseHeaderCondition**  <br />*  **ResponseStatusCondition**  <br />*  **UserAgentCondition**  <br />*  **AsNumberCondition**<br />This value is required when using `additionalMatchOptions`. |
| `positiveMatch` | boolean      | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `additionalMatchOptions`. |
| `values`        | string       | List of values to match on. This value is required when using `additionalMatchOptions`. |


**The apiSelectors object**

Specifies the API endpoints to match on. Note that this object can only be used if the `matchType` is set to **api**. For example:

```
"apiSelectors": [
 {
  "apiDefinitionId": 602,
  "resourceIds": [
   748
   ],
  "undefinedResources": false,
  "definedResources": false
 }
]
```

 Properties of the `apiSelectors` object are described in the following table:

| **Property**         | **Datatype**     | **Description**                                              |
| -------------------- | ---------------- | ------------------------------------------------------------ |
| `apiDefinitionId`    | integer          | Unique identifier of the API endpoint. This value is required when using `apiSelectors`. |
| `resourceIds`        | array  (integer) | Unique identifiers of one or more API endpoint resources.   |
| `undefinedResources` | boolean          | If **true**, matches any resource not explicitly added to your API definition (but without having to include the resource ID in the `resourceIds`  property) . If **false**, matches only those undefined resources listed in the `resourceIds`  property. |
| `definedResources`   | boolean          | If **true**, matches any resource explicitly added to your API definition (but without having to include the resource ID in the `resourceIds` property). If **false**, matches only those defined resources listed in the `resourceIds` property. |


**The bodyParameters object**

Specifies the request body parameters to match on. For example:

```
"bodyParameters": [
 {
  "name": "Country",
  "values": [
   "US",
   "MX",
   "CA"
   ],
  "positiveMatch": true,
  "valueInRange": false
  }
 ]
```

 Properties for the `bodyParameters` object are described in the following table:

| **Property**    | **Datatype** | **Description**                                              |
| --------------- | ------------ | ------------------------------------------------------------ |
| `name`          | string       | Name of the body parameter to match on. This value is required when using `bodyParameters`. |
| `positiveMatch` | boolean      | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `bodyParameters`. |
| `valueInRange`  | boolean      | When **true**, matches values inside the `values` range. Note that your values must be specified as a range to use this property. For example, if your value range is **2:6**, any value between 2 and 6 (inclusive) is a match; values such as 1, 7, 9, or 14 do not match.  <br />When **false**. matches values that fall outside the specified range. |
| `values`        | string       | Body parameter values to match on. This value is required when using `bodyParameters`. |



**The fileExtensions object**

Specifies the file extensions to match on. For example:

```
"fileExtensions": {
 "positiveMatch": false,
 "values": [
  "avi",
  "bmp",
  "jpg"
  ]
 }
```

 Properties of the `fileExtensions` object are described in the following table:

| **Property**    | **Datatype** | **Description**                                              |
| --------------- | ------------ | ------------------------------------------------------------ |
| `positiveMatch` | boolean      | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `fileExtensions`. |
| `values`        | string       | List of file extensions to match on. This value is required when using `fileExtensions`. |


**The path object**

Specifies the paths to match on. For example:

```
"path": {
 "positiveMatch": true,
 "values": [
  "/login/",
  "/user/"
 ]
}
```

Properties of the `path` object are described in the following table:

| **Property**    | **Datatype** | **Description**                                              |
| --------------- | ------------ | ------------------------------------------------------------ |
| `positiveMatch` | boolean      | Indicates whether the policy is triggered if a match is found (**true**) or if a match is  not found (**false**). This value is required when using `path`. |
| `values`        | array        | List of paths to match on. This value is required when using `path`. |


**The queryParameters object**

Specifies the query parameters to match on. For example:

```
"queryParameters": [
 {
  "name": "productId",
  "values": [
  "DOC_12",
  "DOC_11"
  ],
 "positiveMatch": true,
 "valueInRange": false
 }
]
```

Properties of the `queryParameters` object are described in the following table:

| **Property**    | **Datatype** | **Description**                                              |
| --------------- | ------------ | ------------------------------------------------------------ |
| `name`          | string       | Name of the query parameter to match on.  his value is required when using `queryParameters`. |
| `positiveMatch` | boolean      | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not **found** (false). This value is required when using `queryParameters`. |
| `valueInRange`  | boolean      | When **true**, matches values inside the `values` range. Note that your values must be specified as a range to use this property. For example, if your value range is **2:6**, any value between 2 and 6 (inclusive) is a match; values such as 1, 7, 9, or 14 do not match.     <br /><br />When **false**. matches values that fall outside the specified range. |
| `values`        | string       | List of query parameter values to match on. This value is required when using `queryParameters`. |



------

## <a id="match"></a>Create a match target

[Back to table of contents](#contents)

[The match target JSON file](#matchjson)
[Match target sequencing](#sequence)
[Determining the correct match target sequence](#correct)
[The match target sequence JSON file](#sequencejson)

Suppose you hire a clown to perform at your child's birthday party, and the performer promises to come to your office and collect his fee. Because you've just been called into a meeting, you ask a co-worker to give the performer the check when he arrives. Does that mean that your co-worker hands out a check to everyone who walks into the office? Of course not. Instead, you give the assistant very specific instructions: give the check to a man, in his late 40s, who's dressed in a clown suit, and who identifies himself as Mr. Freckles. That's the “match target” for the payment.

Match targets in application security might be a little less colorful, but they serve a somewhat similar purpose. Your website might get millions of requests each day, and you might have any number of security policies that help you handle those requests and help protect your site from potentially malicious requests. Do you apply every single security policy to every single request? No. Instead, you use match targets to define which security policy (if any) should apply to a specific API, hostname, or path. Should a request come in that triggers a match target (for example, you might have a match target that scans for a specific set of file extensions), the security policy associated with the target goes into action, using protections such as rate controls, slow POST protections, and reputation controls to determine whether the request should be honored.

You can create match targets in Terraform by using a configuration similar to this:

```
terraform {
 required_providers {
  akamai  = {
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

resource "akamai_appsec_match_target" "match_target" {
 config_id    = data.akamai_appsec_configuration.configuration.config_id
 match_target = file("${path.module}/match_targets.json")
}
```

In this configuration, we begin by defining **akamai** as our Terraform provider and by providing our authentication credentials. We then use this block to connect to the **Documentation** configuration:

```
data "akamai_appsec_configuration" "configuration" {
 name = "Documentation"
}
```

After the connection is made, we use the [akamai_appsec_match_target](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_match_target) resource and the following block to create the match target:

```
resource "akamai_appsec_match_target" "match_target" {
 config_id    = data.akamai_appsec_configuration.configuration.config_id
 match_target = file("${path.module}/match_targets.json")
}
```

Only two things happen in the preceding block: we specify the ID of the configuration we want the new match target associated with (`config_id`) and we specify the path to the JSON file containing the match target properties and property settings. That's what this argument is for:

```
match_target = file("${path.module}/match_targets.json"
```

In this example, we use a JSON file named **match_targets.json**. That file name is arbitrary: it doesn't have to be named match_targets.json. Similarly, we put our JSON file in the same folder as the Terraform executable; that's what the syntax **${path.module}/** is for. This file path is also arbitrary: you can store your JSON file in any folder you want. Just be sure to replace **\${path.module}/** with the path to that folder.



### <a id="matchjson"></a>The match target JSON file

When you create a match target, the properties and property values for that target are typically defined in a JSON file; when you run your Terraform configuration, information is extracted from that file and used to configure the new match target. A JSON file for creating a match target looks similar to this:

```
{
  "type": "website",
  "isNegativePathMatch": false,
  "isNegativeFileExtensionMatch": false,
   "hostnames": [
    "akamai.com",
    "learn.akamai.com",
    "developer.akamai.com"
  ],
  "fileExtensions": ["sfx", "py", "js", "jar", "html", "exe", "dll", "bat"],
  "securityPolicy": {
    "policyId": "gms1_134637"
  }
}

```

Keep in mind that your match target JSON files won't necessarily look exactly like the preceding file; that's because different match targets have different sets of properties and property values. The following sections of this document provide information on the properties available for use when creating a match target.

#### Required arguments

The following argument must be included in all your match target JSON files:

| **Argument** | **Datatype** | **Description**                                              |
| ------------ | ------------ | ------------------------------------------------------------ |
| `type`       | string       | Match target type. Allowed values are:     <br />*  **website**  <br />*  **api** |


**The securityPolicy object**

Specifies the security policy to be associated with the match target; this object is required in any match target JSON file you create. For example:

```
"securityPolicy": {
 "policyId": "gms1_134637"
 }
```

Arguments related to the `securityPolicy` object are described in the following table:

| **Argument** | **Datatype** | **Description**                            |
| ------------ | ------------ | ------------------------------------------ |
| `policyId`   | string       | Unique identifier of the security policy. |


**Optional arguments**

The arguments described in the following table are optional: they might (or might not) be required depending on the other arguments you include in your match target. For example, if your match target includes the `filePaths` or `fileExtensions` object then your JSON file *can't* include the `defaultFile` argument.

| **Argument**                   | **Datatype** | **Description**                                              |
| ------------------------------ | ------------ | ------------------------------------------------------------ |
| `configId`                     | integer      | Unique identifier of the security configuration containing the match target. |
| `configVersion`                | integer      | Version number of the security configuration associated with the match target. |
| `defaultFile`                  | string       | Specifies how path matching takes place. Allowed values are:     <br />*  **NO_MATCH**. Excludes the default file from path matching.  <br />*  **BASE_MATCH**. Matches only requests for top-level hostnames that end in a  trailing slash.  <br />*  **RECURSIVE_MATCH**. Matches all requests for paths that end in a trailing slash. |
| `fileExtensions`               | array        | File extensions that the match target scans for.            |
| `filePaths`                    | array        | File paths that the match target scans for.                 |
| `hostnames`                    | array        | Hostnames that the match target scans for.                  |
| `isNegativeFileExtensionMatch` | boolean      | If **true**, the match target is triggered if a match *isn't* found in the list of file extensions. |
| `isNegativePathMatch`          | boolean      | If **true**, the match target is triggered if a match *isn't* found in the list of file paths. |
| `sequence`                     | integer      | Ordinal position of the match target in the sequence of match targets. Match targets are processed in the specified order: the match target with the sequence value 1 is processed first, the match target with the sequence value 2 is processed second, etc. |


**The apis object**

Specifies the API endpoints to match on. Note that argument can only be used if the match target's `type` is set to **api**.

Arguments associated with the `apis` object are described in the following table:

| **Argument** | **Datatype** | **Description**                         |
| ------------ | ------------ | --------------------------------------- |
| `id`         | integer      | Unique identifier of the API endpoint. |
| `name`       | string       | Name of the API endpoint name.         |


**The byPassNetworkLists object**

The bypass network list provides a way for you to exempt one or more network lists from the Web Application Firewall. For example:

```
"bypassNetworkLists": [
 {
  "id": "1410_DOCUMENTATIONNETWORK",
  "name": "Documentation Network"
  }
]
```

Arguments associated with the `bypassNetworkLists` object are described in the following table:

| **Argument** | **Datatype** | **Description**                         |
| ------------ | ------------ | --------------------------------------- |
| `id`         | string       | Unique identifier of the network list. |
| `name`       | string       | Name of the network list.              |



### <a id="sequence"></a>Match target sequencing

By default, match targets are applied in the order in which they are created. For example, suppose you have two different match targets:

- One checks to see if there are any file extensions included in a list of file extensions.
- One that checks to see if the hostname is included in the specified list of hostnames.

Assuming that your match targets were created in the order shown above, then for each request:

1. The request is examined to see if it includes any of the specified file extensions.
2. If any of these file extensions are found, the request is examined to see if the hostname is on the list of hostnames.

Although that approach works, it might not be the most efficient route you can take. For example, suppose you have 3 hostnames on your list: hostnames A, B, and C. Let’s further suppose that you get 1 million requests each day. That means that, 1 million times a day, you're checking a request for any of your specified file extensions and, if found, then checking to see if the hostname is on the list of hostnames. That's fine if the majority of your requests come from hostnames A, B, and C. But what if only a handful of requests come from those hosts? That means you're doing a detailed search for file extensions on 1 million requests, even though only 1,000 of those requests are coming from a targeted host.

In a case like that, you might be better off swapping the match target sequence order: start by looking at the hostname on each request instead of starting with the file extensions. After all, if the hostname isn't A, B, or C you're done: there's no need to check the file extensions associated with the request. Instead of checking file extensions on 1 million files, you're checking file extensions only on the 1,000 requests coming from hosts A, B, or C.

> **Note**. A good rule of thumb is to start by applying your most general match targets , and then work down to the more specific match targets. Is the target shape blue? If not, then it doesn't matter. If so, then start to whittle down to questions like is it a blue circle; is it a blue circle with white polka dots; are the polka dots less than 1” in diameter; are those small polka dots oval-shaped rather than circular; etc.

As we learned a moment ago, match targets are applied, by default, in the order in which they were created: the first match target you create is applied first, the second match target you create is applied second, and so on. So what if you want to change the order in which your match targets are applied? You can do that by running a Terraform configuration similar to this:

```
terraform {
 required_providers {
  akamai  = {
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

resource "akamai_appsec_match_target_sequence" "match_targets" {
 config_id             = data.akamai_appsec_configuration.configuration.config_id
 match_target_sequence = file("${path.module}/match_targets_sequence.json")
}
```

All we do here is connect to the **Documentation** security configuration and then call the [akamai_appsec_match_target_sequence](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_match_target_sequence) resource. More specifically, we tell this resource to set the match target order using information found in the file **match_targets_sequence.json**. That's what this argument does:

```
match_target_sequence = file("${path.module}/match_targets_sequence.json")
```

In that line, the value **file("\${path.module}/match_targets_sequence.json")** represents the path to the JSON file that contains the target ordering information. The syntax **${path.module}/** indicates that the JSON file (match_targets_sequence.json) is in the same folder as our Terraform executable file.

>  **Note**. What if you want to put the JSON file in a different folder? That's fine: just be sure you include the full path to that folder in your Terraform configuration.

As you might have guessed, the key to reordering your match targets is to specify the desired target sequence in this JSON file.



### <a id="correct"></a>Determine the current match target sequence

For better or worse, Terraform doesn't return information about match target sequencing; the [akamai_appsec_match_targets](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_match_targets) data source returns only minimal information about your match targets:

```
+---------------------------------+
| matchTargetDS          |
+---------+-------------+---------+
| ID   | POLICYID  | TYPE  |
+---------+-------------+---------+
| 3723387 | gms1_134637 | Website |
| 3722423 | gms1_134637 | Website |
| 3722616 | gms1_134637 | Website |
| 3722692 | gms1_134637 | Website |
| 3723385 | gms1_134637 | Website |
| 3722626 | gms1_134637 | Website |
| 3722379 | gms1_134637 | Website |
+---------+-------------+---------+
```

However, you can use the [Application Security API](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsequence) to return information about your match targets that includes a target's sequence number:

```
"securityPolicy": {
  "policyId": "gms1_134637"
  },
"sequence": 2,
"targetId": 3722423
```



#### <a id="sequencejson"></a>The match target sequence JSON file

The JSON file used to sequence your match targets looks similar to this:

```
{
 "type": "website",
 "targetSequence": [
  {
   "targetId": 3722423,
   "sequence": 1
  },
  {
   "targetId": 2660693,
   "sequence": 2
  },
  {
   "targetId": 2712938,
   "sequence": 3
  },
  {
   "targetId": 2809154,
   "sequence": 4
  },

  {
   "targetId": 3023865,
   "sequence": 5
  },
  {
   "targetId": 3505726,
   "sequence": 6
  },
  {
   "targetId": 3722379,
   "sequence": 7
  }
 ]
}
```

This JSON file has two required properties: `type` (which specifies whether the sequencing is for website matches or API matches), and `targetSequence`, an object containing the `targetId` and `sequence` value for each of your match targets. Do you want match target **3722423** to be the first match target applied? Then set its sequence value to **1**:

```
{
 "targetId": 3722423,
 "sequence": 1
},
```

Continue in this fashion until you've configured all your match targets in the desired order.



------

## <a id="kona"></a>Modify a Kona rule set rule action

[Back to table of contents](#contents)

[View the actions currently assigned to a KRS rule](#view)
[Modify a rule action](#modify)
[Work with custom denies](#denies)

Kona Site Defender uses a vast collection of common vulnerability and exposure (CVE) rules to help protect your website from specific attacks. Each of these rules (collectively referred to as the Kona Rule Set or KRS) is designed to look for a specific exploit and to take action (issue an alert, deny the request, take a custom course of action, or do nothing at all) anytime the rule is triggered. These rule actions are predefined by Akamai, but you can use Terraform to change the action assigned to any of your KRS rules. Do you feel that issuing an alert is not sufficient for a given set of circumstances? Would you prefer that requests be denied any time a specific rule is triggered? Then use Terraform to change the rule action from alert to deny.

In this documentation, we'll show you how to do just that.

>  **Important.** If you’re running Adaptive Security Engine (ASE) in auto mode then you shouldn’t see rule actions; that’s because those actions are automatically configured for you. However, you can still set rule conditions and exceptions when using ASE in auto mode.



### <a id="view"></a>View the action currently assigned to a KRS rule

To view the action currently assigned to a rule, use the [akamai_appsec_rules](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_rules) data source, being sure to specify the ID of the security policy and the ID of the rule you're interested in. For example, this simple API call returns the action assigned to rule **970002** and security policy **gms1_134637**:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
 }

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
 name = var.security_configuration
}
data "akamai_appsec_rules" "rule" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637”
 rule_id            = 970002
}
output "rule_action" {
 value = data.akamai_appsec_rules.rule.rule_action
}
```

That returns information similar to the following:

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

rule_action = "none"
```



### <a id="modify"></a>Modify a rule action

To change the action assigned to a KSD rule, use the [akamai_appsec_rule resource](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rule) and set the rule action to one of the following values:

- **alert**. Writes an entry to the log file any time a request triggers the rule.
- **deny**. Blocks the request using a predefined response.
- **deny_custom_{custom_deny_id}**. Blocks the request using a custom deny response that you create. Custom deny actions are discussed later in this documentation.
- **none**. Takes no action.

For example, the following Terraform configuration sets the rule action for the rule **970002** to **alert**:

```
terraform {
 required_providers {
  akamai  = {
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

resource "akamai_appsec_rule" "rule" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637”
 rule_id            = 970002
 rule_action        = "alert"
 }
```

There's really nothing complicated about this configuration. It begins, like most of our Terraform configurations, by calling the Akamai provider and providing our authentication credentials. After connecting to the **Documentation** security configuration, we then encounter this block:

```
resource "akamai_appsec_rule" "rule" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637”
 rule_id            = 970002
 rule_action        = "alert"
 }
```

Here we use the [akamai_appsec_rule](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rule) resource to change the rule action for our KRS rule. Which KRS rule? That's easy; the rule that:

- Resides in the security configuration we connected to.
- Is associated with the security policy **gms1_134637**.
- Has the rule ID **970002**.

When we run our configuration, we get back output similar to this:

```
akamai_appsec_rule.rule: Creating...
akamai_appsec_rule.rule: Creation complete after 4s [id=90013:gms1_134637:970002]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

And if we rerun our original API call, we should see that the rule action has been changed to alert.

That's all it takes.



### <a id="denies"></a>Work with custom denies

Custom denies provide a way for you to create a custom page or custom API response for rejected requests. These custom pages/responses serve at least two purposes:

- They help you maintain a positive and branded experience in case of a false positive result (e.g., the suspected web attack wasn't actually a web attack).
- They can misdirect actual attackers away from your website.

We won't explain how to create custom denies here; see the documentation for the [akamai_appsec_custom_deny](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_custom_deny) resource for that information. Here, we'll simply show you how to retrieve a collection of your available custom denies, then show you how to use one of those denies as a rule action.



#### View the custom denies available for use

To determine the custom denies available for use, we can use a Terraform configuration similar to this one:

```
terraform {
 required_providers {
  akamai  = {
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

data "akamai_appsec_custom_deny" "custom_deny_list" {
 config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "custom_deny_list_output" {
 value = data.akamai_appsec_custom_deny.custom_deny_list.output_text
}
```

Again, there's nothing very complicated here: we connect to the **Documentation** security configuration, use the [akamai_appsec_custom_deny](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_custom_deny) data source to retrieve a collection of custom denies, then echo back the contents of that collection. As you can see, the **Documentation** configuration has a pair of custom denies:

```
+-------------------------------------+
| customDenyDS            |
+-------------------+-----------------+
| ID        | NAME      |
+-------------------+-----------------+
| deny_custom_64386 | Operation    |
| deny_custom_68193 | new custom deny |
+-------------------+-----------------+
```

If we'd like detailed information about any one of these custom denies, we can essentiaslly rerun this same configuration; the only thing we do different is reference the ID of the custom deny of interest (in this example, **deny_custom_64386**):

```
data "akamai_appsec_custom_deny" "custom_deny_list" {
 config_id      = data.akamai_appsec_configuration.configuration.config_id
 custom_deny_id = "deny_custom_64386"
}
```

If we now run `terraform plan` that should tell us everything we want to know about the custom deny:

```
{
 customDenyList = [
   {
     description = "Operation"
     id     = "deny_custom_64386"
     name    = "Operation"
     parameters = [
       {
         name = "prevent_browser_cache"
         value = "true"
        },
       {
         name = "response_body_content"
         value = <<-EOT
            %(AK_REFERENCE_ID)
            <h1>
            This is my custom Error Message
            </h1>
          EOT
        },
       {
         name = "response_content_type"
         value = "text/html"
        },
       {
         name = "response_status_code"
         value = "403"
        },
      ]
    }
   ]
  }
 )
```

To set a rule action to this custom deny we use the exact same Terraform configuration used when setting the rule action to alert. The one difference? Now we set the `rule_action` argument to the ID of our custom deny:

```
resource "akamai_appsec_rule" "rule" {
 config_id          = data.akamai_appsec_configuration.configuration.config_id
 security_policy_id = "gms1_134637"
 rule_id            = 970002
 rule_action        = "deny_custom_64386"
 }
```



------

## <a id="import"></a>Import a Terraform resource from one security configuration to another

[Back to table of contents](#contents)

[Before we begin](#before)
[Export data from a security configuration](#export)
[Modify the export file](#modifyexport)
[Import the exported settings](#import)

If you have multiple security configurations, it's possible that you have elements (security policies, configuration settings, rate policies, etc.) in configuration A that you'd also like to have in configuration B. For example, you might want the same set of custom rules in both configurations. So how can you replicate items found in configuration A to configuration B? Well, obviously you can simply recreate those items, from scratch, in configuration B: create new security policies with the same exact values, create new configuration settings with the same exact values, create new rate policies with the same exact values. That works, but that process can also be slow, tedious, and prone to errors.

If you're thinking, “There must be a better, easier way to do this,” then give yourself a pat on the back: there *is* a better, easier way to do this. As it turns out, Terraform provides a way for you to export data from configuration A and then import that data into configuration B. The process isn't fully automated – you have to do some manual editing here and there – but the approach is still faster and easier than manually creating new items in configuration B.

In this article we show a simple example: we export the prefetch settings from configuration A (or, more correctly, the configuration with the ID **90013**) and then import those same settings into configuration B (i.e., the configuration with the ID **11118**). We'll focus on the prefetch settings simply, so we can keep our sample files short, sweet, and easy to follow. However, when it comes to exporting and importing items from a configuration you aren't limited just to prefetch settings; instead, you can export and import any (or all) of the following:

- AdvancedSettingsPrefetch
- ApiRequestConstraints
- AttackGroup
- AttackGroupConditionException
- CustomDeny
- CustomRule
- CustomRuleAction
- Eval
- EvalRuleConditionException
- IPGeoFirewall
- MatchTarget
- PenaltyBox
- RatePolicy
- RatePolicyAction
- ReputationProfile
- ReputationProfileAction
- Rule
- RuleConditionException
- SecurityPolicy
- SiemSettings
- SlowPost



#### A Note About the Terraform import Command

In this documentation, we'll use the [akamai_appsec_export_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_export_configuration) data source to retrieve the prefetch settings from configuration **90013** and save those values to a local file (**export.tf**). We'll then add some additional code to export.tf, and use Terraform's `apply` command to assign those settings to configuration **11118**.

If you're familiar with the Terraform language, you might be aware of the `import` command, a command that enables you to import items from one configuration to another. So why don't we just use the `import` command to carry out this operation? Well, that's mainly because the `import` command doesn't actually import items from configuration A into configuration B. Instead, it imports items from configuration A into your Terraform state file (e.g., **terraform.tfstate**). Because the imported data still needs some manual modification, you're left with two choices: directly manipulate the state file (not recommended), or copy the information from the Terraform state to a standalone configuration file. That results in some additional work, and some additional chances for errors to occur.

Which is why we use the [akamai_appsec_export_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_export_configuration) data source instead.

### <a id="before"></a>Before we begin

We'll explain how the export/import process works momentarily. Before we do that, however, we should clarify exactly what that process is going to do. As noted, we have two security configurations: **90013** and **11118**. Currently, the prefetch settings for the two configurations looks like this

| **Property**           | **90013** | **11118** |
| ---------------------- | --------- | --------- |
| `enable_app_layer`     | true      | false     |
| `all_extensions`       | false     | false     |
| `enable_rate_controls` | true      | false     |
| `extensions`           | mp4       |           |

What we want to do, in effect, is copy the setting values from **90013** and use those values to configure the settings for **11118**. If we succeed, the two configurations will have the exact same prefetch setting values. In other words:

| **Property**           | **90013** | **11118** |
| ---------------------- | --------- | --------- |
| `enable_app_layer`     | true      | true      |
| `all_extensions`       | false     | false     |
| `enable_rate_controls` | true      | true      |
| `extensions`           | mp4       | mp4       |

And now we're ready to talk about how you do this.



### <a id="export"></a>Export data from a security configuration

To export data from a security configuration, use the [akamai_appsec_export_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_export_configuration) data source, taking care to specify exactly what it is you want to export. For example, a Terraform configuration that exports prefetch settings looks similar to this:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
 name = "Configuration A"
}

data "akamai_appsec_export_configuration" "export" {
 config_id = data.akamai_appsec_configuration.configuration.config_id
 version   = 9
 search    = ["AdvancedSettingsPrefetch.tf"]
}

resource "local_file" "config" {
 filename = "${path.module}/export.tf"
 content  = data.akamai_appsec_export_configuration.export.output_text
}
```

As usual, there's nothing particularly special about the configuration. It begins, as all our configurations do, by calling the Akamai Terraform provider and by providing our authentication credentials. The configuration then uses this block to connect to **Configuration A**, the security configuration that has the setting values we want to export:

```
data "akamai_appsec_configuration" "configuration" {
 name = "Configuration A"
}
```

The exporting itself takes place with this block:

```
data "akamai_appsec_export_configuration" "export" {
 config_id = data.akamai_appsec_configuration.configuration.config_id
 version   = 9
 search    = ["AdvancedSettingsPrefetch"]
}
```

Here, we're simply calling the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_configuration) data source and connecting that data source to our security configuration (technically, to **version 9** of our security configuration). We then use the `search` argument to specify the items we want to export; in this case that's only the prefetch setting values:

```
search = ["AdvancedSettingsPrefetch"]
```

What if we have additional items we want to export? That's fine: we can just add those items to the search list:

```
search = ["AdvancedSettingsPrefetch", "MatchTarget", "SIEMSettings"]
```

When that's done, we then use this block to tell Terraform what to do with the exported values:

```
resource "local_file" "config" {
 filename = "${path.module}/export.tf"
 content  = data.akamai_appsec_export_configuration.export.output_text
}
```

And what are we doing with the exported values? As it turns out, we're doing two things here. First, we're creating a local file named **export.tf**:

```
filename = "${path.module}/export.tf"
```

>  **Note**. The syntax **\$(path.module)/** indicates that we want to create the file in the same folder as the Terraform executable. But that's entirely up to you: if you'd rather save the file to a different folder, just replace **\$(path.module)/** with the path to that folder.

After that, we tell Terraform to take all the exported settings and setting values and write them to this new file:

```
content = data.akamai_appsec_export_configuration.export.output_text
```

Believe it or not, that's all we have to do. When we run our configuration, we'll get a new file (export.tf) that looks something like this:

```
// terraform import akamai_appsec_advanced_settings_prefetch.akamai_appsec_advanced_settings_prefetch 90013
resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 90013
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```

This file consists of two parts. Part 1 is a commented-out Terraform import command (in a Terraform configuration file, **//** is used to indicate a comment):

```
// terraform import akamai_appsec_advanced_settings_prefetch.akamai_appsec_advanced_settings_prefetch 90013
```

Part 2 consists of the settings exported from configuration **90013**:

```
resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 90013
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```

These are the values we want copied to the target configuration. Before we can do that, however, we need to manually make a change to export.tf.


### <a id="modifyexport"></a>Modify the export file

When you export data from configuration **90013**, the resulting data file is – not surprisingly – all about configuration **90013**. That's why you'll occasionally see references to the configuration ID in the file:

```
config_id = 90013
```

Like we said, that's to be expected: after all, these are the settings and values for configuration **90013**. However, these are also the settings and values we're going to import into a different security configuration. If we leave the file as-is, those settings and values will be imported into configuration **90013**; in other words, we'd import them right back into the configuration we just exported them from. Needless to say, that's not what we want; instead, we want to import those settings into configuration **11118**. That means we need to change all references of **90013** to **11118**. For example:

```
config_id = 11118
```

In our sample file, there's only one place where we need to make a change. However, if we exported multiple items we'll likely have to make this change (or a similar change) in multiple places.

After updating the `config_id` value our export file looks like this:

```
// terraform import akamai_appsec_advanced_settings_prefetch.akamai_appsec_advanced_settings_prefetch 90013

resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 11118
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```

Because the `export` command doesn't create blocks for declaring the Akamai provider and for presenting your Akamai credentials, you'll also need to add this information to the beginning of the export.tf file:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}
```

Oh, and what about the commented-out `import` command. Well, like we said, it's commented out which means you can leave it in or take it out: it won't affect the Terraform configuration in any way. But if you prefer your Terraform configurations be as short and sweet as possible, then you can delete the comment, meaning that the final version of export.tf looks like this:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 11118
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```



### <a id="import"></a>Import the exported settings

We now have a Terraform configuration containing everything we need to apply the same prefetch settings found in configuration **90013** to configuration **11118**. That means that we can use the `terraform plan` command to do a quick syntax check, and then use `terraform apply` to import settings. (Just like we'd do if we were updating the settings from scratch, without copying the values found in configuration **90013**.)

After running `terraform apply`, we use the [akamai_appsec_advanced_settings_prefetch](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_advanced_settings_prefetch) data source to verify that our setting values have been updated:

```
+----------------------------------------------------------------------+
| advancedSettingsPrefetchDS                      |
+------------------+---------------+----------------------+------------+
| ENABLE APP LAYER | ALL EXTENSION | ENABLE RATE CONTROLS | EXTENSIONS |
+------------------+---------------+----------------------+------------+
| true       | false     | true         | mp4    |
+------------------+---------------+----------------------+------------+
```



#### One thing to watch for when working with multiple .tf files

When working in a folder that contains multiple .tf files, (in our case, akamai.tf and export.tf) you might run into problems similar to this:

```
│ Error: Duplicate required providers configuration
│
│  on export.tf line 2, in terraform:
│  2:  required_providers {
│
│ A module may have only one required providers configuration. The required providers were previously configured at akamai.tf:2,3-21.
```

In this case, the “duplicate providers” error occurs because we've used this block in two different files (i.e., in both akamai.tf and export.tf):

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}
```

Does it matter that we reference the Akamai provider in two different .tf files? Well, if those two files reside in the same folder, yes, it does. That's because of the way Terraform processes .tf files. By default, when you run a Terraform command like `terraform plan` or `terraform apply`, Terraform runs all the .tf files it finds in the working folder. That's why you never have to tell Terraform which .tf file to run: it's going to run all of them.

However, Terraform isn't going to run these files one-by-one; that is, it isn't going to run akamai.tf and then run export.tf. Instead, it's effectively going to combine both of those files into a single file, and then run that combined file. And that's where the problem occurs: both akamai.tf and export.tf have blocks that call the Akamai provider, and you can't call the same provider twice in a single configuration. The net result? The “duplicate providers” errors.

So how do you get around this issue?

One obvious way is to save export.tf to a different folder; if you do that, then you won't have to worry about having multiple .tf files in the same folder. Another solution is to temporarily rename the file akamai.tf; for example, you might tack a new file extension on the end (e.g., **.temp**) meaning you'll now have these two files:

- akamai.tf.temp
- export.tf

With akamai.tf temporarily renamed, you have just one .tf file left. Problem solved.

Alternatively, you can comment out everything in akamai.tf; as you probably know, commented lines aren't executed, which means those lines won't conflict with the content in export. tf. You can quickly comment out an entire .tf file by making **/*** the first line in the file and ***/** the last line in the file; everything in between will be commented out. For example:

```
/*
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}
*/
```

Depending on the editor you used to create your Terraform configurations, you'll also have a visual cue telling you that the lines won't be executed.

To restore functionality, remove the two comment markers (**/*** and ***/**).



------

## <a id="aag"></a>Create an automated attack groups (AAG) security configuration

[Back to table of contents](#contents)

When Akamai released automated attack groups (AAG) in October 2018, the technology represented a revolutionary development in protecting websites against Internet attacks. Prior to AAG, attack management primarily revolved around rule management: downloading rules to protect websites from common vulnerabilities and exposures (CVE), configuring rules, updating rules, deleting obsolete – you get the idea. And this approach worked: websites were better protected. The downside? At times, management if these rules could be a little difficult and a little time-consuming.

Automated attack groups offer a much-more scalable, cloud-based approach to protecting websites. Rules still play an important role in AAG, but administrators don’t have to manage those rules one-by-one. Instead, rules are divided into a set of attack groups (for example, one such group helps protect sites from SQL injection attacks), and administrators only have to decide which attack groups they want to deploy. Likewise, there's no need to upload new rules or delete obsolete rules: instead, the rules included in an attack group are managed (e.g., uploaded and deleted) by Akamai.

This article explains how you can use Terraform to create a security configuration that leverages automated attack groups. To do this requires a Terraform configuration that can carry out the following steps:

1. [Create an IP network list](#step1)
2. [Activate the network list](#step2)
3. [Create a security configuration](#step3)
4. [Create a security policy](#step4)
5. [Assign an IP network list](#step5)
6. [Enable network protections](#step6)
7. [Create and configure rate policies](#step7)
8. [Enable rate control protections](#step8)
9. [Configure Logging Settings](#step9)
10. [Configure prefetch settings](#step10)
11. [Enable slow POST protections](#step11)
12. [Configure Slow Post Settings](#step12)
13. [Enable Web Application Firewall protections](#step13)
14. [Configure the Web Application Firewall mode](#step14)
15. [Configure attack group settings](#step15)
16. [Enable and configure the penalty box](#step16)
17. [Activate the security configuration](#step17)

In this article, we explain how to use Terraform to carry out each one of these steps. However, we'll start by showing you a Terraform configuration that carries out all 18 steps.



### The AAG Terraform configuration

The Terraform configuration that creates our AAG security configuration is shown below. Admittedly, the configuration might look a bit intimidating at first; that's mainly because this one configuration is carrying out 17 separate actions. (Or even a few more, depending on what you want to count as a single action.) But don't worry: the configuration's bark is far worse than its bite. And to prove that, we’ll walk you through each of the blocks included in this configuration.

But first, here's a Terraform configuration that creates and enables an AAG security configuration:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

// Step 1: Create an IP Network List
resource "akamai_networklist_network_list" "network_list" {
 name        = "Documentation Test Network"
 type        = "IP"
 description = "Network used for the AAG documentation example."
 list        = ["192.168.1.1","192.168.1.2","192.168.1.3","192.168.1.4"]
 mode        = "REPLACE"
}

// Step 2: Activate the Network List
resource "akamai_networklist_activations" "activation" {
 network_list_id     = akamai_networklist_network_list.network_list.uniqueid
 network             = "Documentation Test Network"
 notes               = "Activation of the AAG test network."
 notification_emails = ["gstemp@akamai.com","karim.nafir@mail.com"]
}

// Step 3: Create a Security Configuration
resource "akamai_appsec_configuration" "create_config" {
 name    = "Documentation AAG Test Configuration"
 description = "This security configuration is used by the documentation team for testing purposes."
 contract_id = "1-3UW382"
 group_id    = 13139
 host_names  = ["llin.gsshappylearning.com"]
}

// Step 4: Create a Security Policy
resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id              = akamai_appsec_configuration.create_config.config_id
 default_settings       = true
 security_policy_name   = "Documentation Security Policy"
 security_policy_prefix = "doc0"
}

// Step 5: Assign an IP Network List
resource "akamai_appsec_ip_geo" "akamai_appsec_ip_geo" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 mode               = "allow"
 ip_network_lists   = [akamai_networklist_network_list.network_list.uniqueid]
}

// Step 6: Enable Network Protections
resource "akamai_appsec_ip_geo_protection" "protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}

// Step 7: Create and Configure Rate Policies
resource "akamai_appsec_rate_policy" "rate_policy_1" {
  config_id   = akamai_appsec_configuration.create_config.config_id
  rate_policy = file("${path.module}/rate_policy_1.json")
}

 resource "akamai_appsec_rate_policy_action" "rate_policy_actions_1" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.rate_policy_1.id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
 }

resource "akamai_appsec_rate_policy" "rate_policy_2" {
  config_id   = akamai_appsec_configuration.create_config.config_id
  rate_policy = file("${path.module}/rate_policy_2.json")
}

 resource "akamai_appsec_rate_policy_action" "rate_policy_actions_2" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.rate_policy_2.id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
 }

 resource "akamai_appsec_rate_policy_action" "rate_policy_actions_3" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.rate_policy_3.id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
 }

// Step 8: Enable Rate Control Protections
resource "akamai_appsec_rate_protection" "protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}

// Step 9: Configure Logging Settings
resource "akamai_appsec_advanced_settings_logging" "logging" {
 config_id = akamai_appsec_configuration.create_config.config_id
 logging   = file("${path.module}/logging.json")
}

// Step 10: Configure Prefetch Settings
resource "akamai_appsec_advanced_settings_prefetch" "prefetch" {
 config_id            = akamai_appsec_configuration.create_config.config_id
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = false
 extensions           = ["cgi","jsp","aspx","EMPTY_STRING","php","py","asp"]
}

// Step 11: Enable Slow Post Protections
resource "akamai_appsec_slowpost_protection" "protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}

// Step 12: Configure Slow Post Settings
resource "akamai_appsec_slow_post" "slow_post" {
 config_id                  = akamai_appsec_configuration.create_config.config_id
 security_policy_id         = akamai_appsec_security_policy.security_policy_create.security_policy_id
 slow_rate_threshold_rate   = 10
 slow_rate_threshold_period = 30
 duration_threshold_timeout = 20
 slow_rate_action           = "alert"
}

// Step 13: Enable Web Application Firewall Protections
resource "akamai_appsec_waf_protection" "akamai_appsec_waf_protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}

// Step 14: Configure the Web Application Firewall Mode
resource "akamai_appsec_waf_mode" "waf_mode" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 mode               = "AAG"
}

// Step 15: Configure Attack Group Settings
resource "akamai_appsec_attack_group" "akamai_appsec_attack_group_AAG1" {
 for_each            = toset(["SQL", "XSS", "CMD", "HTTP", "RFI", "PHP", "TROJAN", "DDOS", "IN", "OUT"])
 config_id           = akamai_appsec_configuration.create_config.config_id
 security_policy_id  = akamai_appsec_security_policy.security_policy_create.security_policy_id
 attack_group        = each.value
 attack_group_action = "deny"
}

// Step 16: Enable and Configure the Penalty Box
resource "akamai_appsec_penalty_box" "penalty_box" {
 config_id              = akamai_appsec_configuration.create_config.config_id
 security_policy_id     = akamai_appsec_security_policy.security_policy_create.security_policy_id
 penalty_box_protection = true
 penalty_box_action     = "alert"
}

// Step 17: Activate the Security Configuration
resource "akamai_appsec_activations" "new_activation" {
  config_id           = akamai_appsec_configuration.create_config.config_id
  network             = "STAGING"
  notes               = "Activates the Documentation AAG Test Configuration on the staging network."
  activate            = true
  notification_emails = ["gstemp@akamai.com","karim.nafir@mail.com"]
}
```

Let's see if we can explain exactly what this configuration does, and how it goes about doing it.



### <a id="step1"></a>Step 1: Create an IP network list

```
resource "akamai_networklist_network_list" "network_list" {
 name        = "Documentation Test Network"
 type        = "IP"
 description = "Network used for the AAG documentation example."
 list        = ["192.168.1.1","192.168.1.2","192.168.1.3","192.168.1.4"]
 mode        = "REPLACE"
}
```

Among other things, network lists provide a way for you (by way of your firewall) to manage clients based either on their IP address or on their geographic location. For example, if you want to prevent  IP address 192.168.1.0 through 192.168.1.255 from passing through your firewall, you set the `list` property of the network list to the Classless Inter-Domain Routing (CIDR) address **192.168.1.0/24**. Alternatively, you can block (or allow) all clients from Norway by setting the geographic list property to **NO**, the ISO 3166 country code for Norway. Or you can do both: Akamai enables you to create as many network lists as you need.

To create a network list, use the [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list) resource and the following arguments:

| **Argument**  | **Description**                                              |
| ------------- | ------------------------------------------------------------ |
| `name`        | Name of the new network list. Network names don't have to be unique: you can have multiple network lists that share the same name. However, when the list is created it’s issued a unique ID, a value comprised of a numeric prefix and the list name. (Or a variation of that name. For example, a list named **Documentation Network** will be given an ID similar to **108970_DOCUMENTATIONNETWORK**, with the blank space in the name being removed.) |
| `type`        | Indicates the type of addresses used on the list. Allowed values are:     <br />*  **IP**. For IP/CIDR addresses.  <br />*  **GEO**. For ISO 3166 geographic codes.     <br /><br />Note that you can’t mix IP/CIDR addresses and geographic codes on the same list. |
| `description` | Brief  description of the network list.                      |
| `list`        | Array containing either the IP/CIDR addresses or the geographic codes to be added  to the new network list. For example:     <br /><br />`list = ["US", "CA", "MX"]`     <br /><br />Note that the list value is formatted as an array even if you only add a single  item to that list:     <br /><br />`list = ["US"]`     <br /><br />Note, too that `list` is the one optional argument available to you: you don't have to include this argument in your configuration. However, leaving out the `list` argument also means that you'll create a network list that has no IP/CIDR addresses or geographic  codes. |
| `mode`        | Set to **REPLACE** when creating a new network list.        |



### <a id="step2"></a>Step 2: Activate the network list

```
resource "akamai_networklist_activations" "activation" {
 network_list_id     = akamai_networklist_network_list.network_list.uniqueid
 network             = "STAGING"
 notes               = "Activation of the AAG test network."
 notification_emails = ["gstemp@akamai.com","karim.nafir@mail.com"]
}
```

After a network list has been created your firewall can’t block (or allow) clients on that list until that list has been activated. You can read more about network activation in the [Network Lists Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_networklists). In the meantime, network lists can be activated by using the [akamai_networklist_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_activations) resource and the following arguments:

| **Argument**          | **Description**                                              |
| --------------------- | ------------------------------------------------------------ |
| `network_list_id`     | Unique  identifier of the network list being activated. In our AAG Terraform configuration we refer to the network list ID like this:     <br /><br />`akamai_networklist_network_list.network_list.uniqueid`     <br /><br />That, as you might recall, is the ID assigned to the network list created in the  previous step. (It probably goes without saying that we can't hardcode a network list ID in our configuration: after all, the ID won't exist until after we've called `terraform apply` and the network list has been created.) |
| `network`             | Specifies the network that the network list is being activated for. Allowed values are:     <br />*  **staging**. “Sandbox” network used for testing and fine-tuning. The staging network includes a small subset of Akamai edge servers but is not used to protect your actual website.  <br />*  **production**. Network lists activated on the production network help protect your actual website.     <br /><br />If this argument is omitted, the network list is automatically activated on the staging network |
| `notes`               | Arbitrary information about the network list and its activation. |
| `notification_emails` | JSON array of email addresses of the people to be notified when the activation process finishes. |



### <a id="step3"></a>Step 3: Create a security configuration

```
resource "akamai_appsec_configuration" "create_config" {
 name    = "Documentation AAG Test Configuration"
 description = "This security configuration is used by the documentation team for testing purposes."
 contract_id = "1-3UW382"
 group_id    = 13139
 host_names  = ["llin.gsshappylearning.com"]
}
```

Security configurations are containers that house all the elements – security policies, rate policies, match targets, slow POST protection settings – that make up a website protection strategy. Before you can create any of these elements you need a place for those elements to reside. That place is a security configuration.

Security configurations are created by using the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_configuration) resource and the following arguments:

| **Argument**  | **Description**                                              |
| ------------- | ------------------------------------------------------------ |
| `name`        | Unique name to be assigned to the new configuration.        |
| `description` | Brief description of the configuration and its intended purpose. |
| `contract_id` | Akamai contract ID associated with the new configuration. You can use the [akamai_appsec_contracts_groups](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_contracts_groups) data source to return information about the contracts and groups available to you. |
| group_id      | Akamai group ID associated with the new configuration.      |
| host_names    | Names of the selectable hosts to be protected by the configuration. Note that names must be passed as an array; that's what the square brackets surrounding **"documentation.akamai.com"** are for. To add multiple hostnames to the configuration, separate the individual names by using commas. For example:     <br /><br />`host_names = ["documentation.akamai.com", "training.akamai.com", "events.akamai.com"]`     <br /><br />All security configurations must include at least one protected host. |

For more information about creating a security configuration, see **Creating Security Configurations**.



### <a id="step4"></a>Step 4: Create a security policy

```
resource "akamai_appsec_security_policy" "security_policy_create" {
 config_id              = akamai_appsec_configuration.create_config.config_id
 default_settings       = false
 security_policy_name   = "Documentation Security Policy"
 security_policy_prefix = "doc0"
}
```

Security policies are probably the single most important item found in a security configuration. That shouldn't come as much of a surprise: after all, many of the other items used in a security configuration (rate policies, attack group settings, firewall allow and block lists, slow POST protection settings, etc., etc.) must be associated with a security policy. Although you can create a security configuration without creating a security policy, a policy-less security configuration is of very little use.

Because of that, as soon as we create our security configuration we add a security policy to the configuration. (A security configuration can contain multiple security policies, but we'll create just one for now.) To create the policy, we'll use the [akamai_appsec_security_policy](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_security_policy) resource and the following arguments:

| **Argument**           | **Description**                                              |
| ---------------------- | ------------------------------------------------------------ |
| config_id              | Unique identifier of the security configuration to be associated with the new  policy. Note that, in our AAG Terraform configuration example we always refer to the security configuration ID like this: <br />`akamai_appsec_configuration.create_config.config_id`     <br /><br />That value represents the ID of the security configuration we just created. Needless to say, we can't hardcode the security configuration ID because that ID doesn’t exist until we run terraform apply and create the configuration. |
| security_policy_name   | Unique name to be assigned to the new policy.               |
| security_policy_prefix | Four-character prefix used to construct the security policy ID. For example, a policy with the ID **gms1_134637** is composed of three parts:     <br /><br />*  The security policy prefix (**gms1**)  <br />*  An underscore **(_**)  <br />*  A random value supplied when the policy is created (**134637**) |
| default_settings       | If **true**, the policy is created using the default settings for a new security policy. If **false**, a “blank” security policy is created instead. In our sample configuration we'll set this value to **false**. |

For more information about creating security policies, see the article Creating a Security Policy.



### <a id="step5"></a>Step 5: Assign an IP network list

```
resource "akamai_appsec_ip_geo" "akamai_appsec_ip_geo" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 mode               = "block"
 ip_network_lists   = [akamai_networklist_network_list.network_list.uniqueid]
}
```

As noted previously, network lists enable you to configure your firewall to automatically block (or allow) a set of clients based on either IP address or geographic location. After your list (or lists) have been created, use the [akamai_appsec_ip_geo](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_ip_geo) resource to specify what you want done with these lists. The akamai_appsec_ip_geo resource accepts the following arguments:

| **Argument**                 | **Description**                                              |
| ---------------------------- | ------------------------------------------------------------ |
| `config_id`                  | Unique identifier of the security configuration associated with the network list. This is a required property. |
| `security_policy_id`         | Unique identifier of the security policy associated with the network lists. This is a required property. |
| `mode`                       | Indicates whether the networks that appear on either the geographic networks list or on the IP network list should be allowed to pass through the firewall. Valid  values are:     <br /><br />*  **allow**. Only networks on the geographic/IP network lists are allowed through the firewall. All other networks are blocked.  <br />*  **block**. All networks are allowed through the firewall except for networks on the geographic/IP network lists. Clients on those networks are blocked.     <br /><br />This is a required property. |
| `geo_network_lists`          | Geographic networks on this list are either allowed or blocked based on the value of the `mode` argument. |
| `ip_network_lists`           | Geographic networks on this list are either allowed or blocked based on the value of the `mode` argument. |
| `exception_ip_network_lists` | Networks on this list are always allowed through the firewall, regardless of the networks that do (or don't) appear on either the geographic or IP networks  list. |

In our sample Terraform block, we use these two lines to block all the network lists associated with the security policy:

```
mode = "block"
ip_network_lists = [akamai_networklist_network_list.network_list.uniqueid]
```

This means that only clients that appear on the IP or the geographic network list are prevented from going through the firewall.



### <a id="step6"></a>Step 6: Enable network protections

```
resource "akamai_appsec_ip_geo_protection" "protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}
```

After you've configured your network lists the next step is to enable network protections; if you don't, your security policy won't enforce any of the settings applied to those lists. (For example, any lists you've designated for blocking won't actually be blocked.) To enable network protections, use the [akamai_appsec_ip_geo_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_ip_geo_protection) resource and set the `enabled` property to **true**.

Keep in mind that network protections are enabled on a security policy-by-security policy basis. If you have multiple security policies you'll need to enable network protections on each one.



### <a id="step7"></a>Step 7: Create rate policies

```
resource "akamai_appsec_rate_policy" "rate_policy_1" {
  config_id   = akamai_appsec_configuration.create_config.config_id
  rate_policy = file("${path.module}/rate_policy_1.json")
}

resource "akamai_appsec_rate_policy_action" "rate_policy_actions_1" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.rate_policy_1.id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
 }
```

Rate policies help you monitor and moderate the number and rate of all the requests you receive; in turn, this helps you prevent your website from being overwhelmed by a sudden deluge of requests (which could be an attack of some kind or just an unexpected surge in legitimate traffic). You create rate policies by using the [akamai_appsec_rate_policy](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_policy) resource and the following arguments:

| **Argument**  | **Description**                                              |
| ------------- | ------------------------------------------------------------ |
| `config_id`   | Unique identifier of the security configuration associated with the rate policy. |
| `rate_policy` | File path to the JSON file containing configuration information for the rate policy. In our sample configuration, **\$(path.module)/** indicates that the JSON file (**rate_policy_1.json**) is stored in the same folder as the Terraform executable. This isn't required: you can store your JSON files anywhere you want. Just make sure to specify the full path so that Terraform can find those files. |

>  **Note**. We won't delve into the JSON files in this article. See Creating a Rate Policy for more information about what one of these JSON files actually looks like.

That's all you need to do to create a rate policy. However, when you create a rate policy the rate policy action is automatically set to a null value; that means that nothing happens any time the policy is triggered. Because we prefer to have requests be denied if they trigger one of our rate policies, we create the rate policy and then use a Terraform block similar to this to set the action for that policy:

```
resource "akamai_appsec_rate_policy_action" "rate_policy_actions_1" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.rate_policy_1.id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
 }
```

As you can see, this policy sets both the IPv4 and the IPv6 actions to **deny**.

Note that, in our sample Terraform configuration, we create 3 rate policies; as a result, we've repeated the Terraform blocks shown at the beginning of this step 3 times. It's possible to use a for_each loop to repeat the same action multiple times in a single block of code; we'll show you an example of that when we configure attack group settings. In this case, however, we took the easy way out and simply used the same Terraform block over and over (changing only the variable names).

Why did we do that? Well, creating a rate policy action requires you to know, in advance, the rate policy ID. When trying to work in a for_each loop that's tricky; it's easy to encounter an error like this:

```
The "for_each" value depends on resource attributes that cannot be determined until apply, so Terraform cannot predict how many instances will be created.
```

To help you avoid that error, we skipped the whole for-loop thing altogether.


### <a id="step8"></a>Step 8: Enable rate control protections

```
resource "akamai_appsec_rate_protection" "protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}
```

After you've created your rate policies and assigned your rate policy actions, you enable rate control protections on your security policy; if you don't do this your rate control policies won't actually get used. Because rate control policies are enabled on the security policy (as opposed to having to enable each and every policy), all you need to do is:

1. Call the [akamai_appsec_rate_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_protection) resource.
2. Specify the ID of your security configuration and your security policy.
3. Set the `enabled` property to **true**.



### <a id="step9"></a>Step 9: Configure logging settings

```
resource "akamai_appsec_advanced_settings_logging" "logging" {
 config_id = akamai_appsec_configuration.create_config.config_id
 logging   = file("${path.module}/logging.json")
}
```

HTTP requests and responses are always accompanied by an HTTP header; the header contains detailed information about the request/response, including such things as the cookies set or returned, the web browser (user agent) involved in the transaction, etc. For example, an HTTP GET request includes a header with values similar to these:

```
200 OK
Access-Control-Allow-Origin: *
Connection: Keep-Alive
Content-Type: text/html; charset=utf-8
Date: Mon, 2 Aug 2021 12:06:00 GMT
Etag: "3987c68d0ba92bbeb8b0f612a9199fghm3a69hh"
Keep-Alive: timeout=10, max=788
Server: Apache
Set-Cookie: documentation-cookie=test; expires= Mon, 2 Aug 2022 12:06:00 GMT
; Max-Age=31449600; Path=/; secure
Transfer-Encoding: chunked
Vary: Cookie, Accept-Encoding
```

The [akamai_appsec_advanced_settings_logging](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_advanced_settings_logging) resource enables you to specify which HTTP headers you want to log and which ones you don't. You configure this information by using a JSON file similar to the following:

```
{
 "allowSampling": true,
 "cookies": {
  "type": "exclude",
  "values": [
   "documentation-cookie=test”
  ]
 },  
 "customHeaders": {
  "type": "all"
  },
 "standardHeaders": {
  "type": "all"
  },
 "override": false
}
```

The arguments used in the JSON file are described in the following table:

| **Argument**      | **Description**                                              |
| ----------------- | ------------------------------------------------------------ |
| `allowSampling`   | Set to **true** to enable HTTP header logging. Set to **false** to disable header logging. |
| `cookies`         | Specifies how cookie headers (i.e., HTTP headers that reference cookies set by the server) are logged. Allowed values are:     <br />*  **all**. All cookie headers are logged.  <br />*  **none**. No cookie headers are logged.  <br />*  **exclude**. All cookie headers except the ones specified by the `type` argument are logged.  <br />*  **only**. Only the cookie headers specified by the `type` argument are allowed.     <br /><br />For example:  <br /><br />`"cookies":  {  "type":  "exclude",  "values":  [  "documentation-cookie=test”  ]  }` |
| `standardHeaders` | Specifies how standard HTTP headers such as User-Agent, Forwarded, and Referer, should be logged. Allowed values are:     <br />*  **all**. All standard headers are logged.  <br />*  **none**. No standard headers are logged.  <br />*  **exclude**. All standard headers except the ones specified by the `type` argument  are logged.  <br />* **only**. Only the standard headers specified by the `type` argument are allowed.     <br /><br />For example:  <br /><br />`"standardHeaders":  {  "type":  "only",  "values":  [  "User  Agent", "Referer”  ]  }` |
| `customHeaders`  | Specifies how, custom headers (i.e., non-standard HTTP headers) should be logged. Allowed values are:     <br /><br />*  **all**. All custom headers are logged.  <br />*  **none**. No custom headers are logged.  <br />*  **exclude**. All custom headers except the ones specified by the `type` argument are logged.  <br />*  **only**. Only the custom headers specified by the `type` argument are allowed     <br /><br />For  example:  <br /><br />`"customHeaders":  {  "type":  "all"  }` |
| override          | If **true**, header data isn’t logged for any security events triggered by settings in the security configuration. |

To log HTTP headers in a security configuration you need to:

1. Create a JSON file containing the logging criteria.
2. Call the [akamai_appsec_advanced_settings_logging](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_advanced_settings_logging) resource, specifying the ID of the security configuration and the path to the JSON file.

In our sample Terraform block, we use this line to indicate the path to the JSON file:

```
logging = file("${path.module}/logging.json")
```

As we've seen elsewhere, the syntax **\${path.module}/** indicates that the JSON file (**logging.json**) can be found in the same folder as the Terraform executable. This isn't a requirement: you can store the JSON file anywhere you want. Just be sure that, in your Terraform configuration, you include the full path to the file.

In our sample Terraform configuration, logging settings are applied to all the security policies to the security configuration. However, by including the optional `security_policy_id` argument we can apply these values to an individual policy. In a case like that, the logging settings applied to the policy take precedence over the logging settings applied to the security configuration.



### <a id="step10"></a>Step 10: Configure prefetch settings

```
resource "akamai_appsec_advanced_settings_prefetch" "prefetch" {
 config_id            = akamai_appsec_configuration.create_config.config_id
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = false
 extensions           = ["cgi","jsp","aspx","EMPTY_STRING","php","py","asp"]
}
```

By default, your Web Application Firewall only inspects external requests: requests that originate outside of your origin servers and Akamai's edge servers. Internal requests – requests between your origin servers and Akamai's edge servers – typically aren't inspected, and typically don't *need* to be inspected. (As a general rule, these “prefetch” requests are safe, and inspecting each one doesn't do much besides slowing down your website.)

However, there might be times when enabling prefetch is useful l(for example, if you're concerned about prefetch-driven amplification attacks). If so, you can enable and configure your prefetch settings by using the [akamai_appsec_advanced_settings_prefetch](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_advanced_settings_prefetch) resource and the following arguments:

| **Argument**         | **Description**                                              |
| -------------------- | ------------------------------------------------------------ |
| config_id            | Unique identifier of the security configuration associated with the prefetch  settings. |
| enable_app_layer     | Set to **true** to enable prefetch request inspection.       |
| all_extensions       | Set to **true** to enable prefetch request inspections on all file extensions included in a request. To limit the file extensions to a specified set, set this value to **false** and then specify the target file extensions by using the `extensions` argument. |
| enable_rate_controls | Set to **true** to enable rate policy checking on prefetch requests. |
| extensions           | Specifies the file extensions that, when included in a request, trigger a prefetch inspection. Note that this argument should only be included when ` is set to **false**. |

Prefetch settings apply to the entire security configuration.



### <a id="step11"></a>Step 11: Enable slow POST protections

```
resource "akamai_appsec_slowpost_protection" "protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}
```

Denial of service (DOS) attacks are attacks in which a website is inundated with a massive barrage of requests, each request sent in rapid-fire succession. DOS attacks are bad but, unfortunately, they aren't the only way to bring down a website: another common attack vector is to slowly (*very* slowly) send a series of requests to a site. Because the requests, and the responses, take so long, the website spends its time waiting for the client to respond instead of spending its time handling requests from new (and legitimate) clients. To help guard against these slow POST attacks, use the [akamai_appsec_slowpost_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_slowpost_protection) resource and set the `enabled` property to **true**.



### <a id="step12"></a>Step 12: Configure slow POST settings

```
resource "akamai_appsec_slow_post" "slow_post" {
 config_id                  = akamai_appsec_configuration.create_config.config_id
 security_policy_id         = akamai_appsec_security_policy.security_policy_create.security_policy_id
 slow_rate_threshold_rate   = 10
 slow_rate_threshold_period = 30
 duration_threshold_timeout = 20
 slow_rate_action           = "alert"
}
```

After slow POST protections have been enabled, you might want to adjust the slow POST configuration settings as well. That’s done by using the [akamai_appsec_slow_post](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_slow_post) resource and the following arguments:

| **Argument**                 | **Description**                                              |
| ---------------------------- | ------------------------------------------------------------ |
| `config_id`                  | Unique identifier of the security configuration associated with the slow POST settings. |
| `security_policy_id`         | Unique identifier of the security policy associated with the slow POST settings. |
| `slow_rate_threshold_rate`   | Specifies the minimum rate (in bytes per second) that a request must achieve to avoid triggering the slow POST policy. The threshold rate represents the average number of bytes received during the slow rate threshold period. |
| `slow_rate_threshold_period` | Time period (in seconds) used to calculate the slow rate threshold rate. |
| `duration_threshold_timeout` | Specifies the maximum length of time (in seconds) that the server waits for the first 8KB of a POST request body to be received. If the duration threshold expires before the request has completed or before the first 8KB have been received then the slow POST policy is triggered.     <br /><br />Note that the duration threshold always takes precedence over the slow rate  threshold. |
| `slow_rate_action`           | Specifies the action taken if the policy is triggered. Allowed values are:     <br />*  **alert**. An alert is issued.  <br />*  **abort**. The request is abandoned. |



### <a id="step13"></a>Step 13: Enable Web Application Firewall protections

```
resource "akamai_appsec_waf_protection" "akamai_appsec_waf_protection" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 enabled            = true
}
```

In order to use the Web Application Firewall (WAF), that firewall must be enabled. To enable firewall protection, use the [akamai_appsec_waf_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_waf_protection) resource and:

1. Connect to the appropriate security configuration.
2. Connect to the appropriate security policy (WAF is enabled/disabled on a security policy-by-security policy basis).
3. Set the `enabled` property to **true**.

After enabling the firewall you'll also want to configure the firewall mode and configure your attack group settings. That’s what we do is Step 14.



### <a id="step14"></a>Step 14: Configure the Web Application Firewall mode

```
resource "akamai_appsec_waf_mode" "waf_mode" {
 config_id          = akamai_appsec_configuration.create_config.config_id
 security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
 mode               = "AAG"
}
```

The Web Application Firewall mode determines the way in which the rules in your Kona Rule Set are updated. When using automated attack groups, this value is set to AAG: that ensures that Akamai takes care of updating the rules as needed. Setting the firewall mode to KRS puts the onus on you: you'll need to periodically, and manually, update the rules by yourself.

To specify the firewall mode, use the [akamai_appsec_waf_mode](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_waf_mode) resource and, when using attack groups, set the `mode` to **AAG**.



### <a id="step15"></a>Step 15: Configure attack group settings

```
resource "akamai_appsec_attack_group" "akamai_appsec_attack_group_AAG1" {
 for_each            = toset(["SQL", "XSS", "CMD", "HTTP", "RFI", "PHP", "TROJAN", "DDOS", "IN", "OUT"])
 config_id           = akamai_appsec_configuration.create_config.config_id
 security_policy_id  = akamai_appsec_security_policy.security_policy_create.security_policy_id
 attack_group        = each.value
 attack_group_action = "deny"
}
```

The Kona Rule Set consists of scores of individual rules similar to this:

| **Rule ID** | **Description**                                              |
| ----------- | ------------------------------------------------------------ |
| 3000000     | A SQL injection attack consists of insertion or "injection" of a SQL query via the input data from the client to the application. A successful SQL injection exploit can read sensitive data from the database, modify database data (Insert/Update/Delete), execute administration operations on the database (such as shutdown the DBMS), recover the content of a given file present on the DBMS file system and in some cases issue commands to the  operating system.     <br /><br />One of the common ways to probe applications for SQL Injection vulnerabilities is to use the 'GROUP BY' and 'ORDER BY' clause. Prior to using these clauses in a SQL statement, the hacker first terminates the current query's context (assuming user input is used in the WHERE clause), which could be either numeric or a string literal, and after the clauses and comments out the rest of the query.     <br /><br />This rule triggers on HTTP requests, which contain SQL Injection probes that use the 'GROUP BY' and 'ORDER BY' clause as mentioned above, when they are sent as user-input. |

These rules have further been classified into the following categories:

- **SQL** (SQL Injection). Attack type in which malicious SQL queries are inserted into a data entry field and then executed. Execution of these queries often results in the attacker gaining access to personally-identifiable information about a website's users.
- **XSS** (Cross-Site Scripting). Attack type in which client-side scripts are added to a web page and thus made available to users who unwittingly execute those scripts.
- **CMD** (Command Injection). Attack type that enables arbitrary (and typically malicious) commands to be executed on a host's operating system.
- **HTTP**. (HTTP Injection). Attack type in which malicious commands are included within the parameters of an HTTP request.
- **RFI**. (Remote File Inclusion). Attack type in which a malefactor attempts to dynamically insert malicious code into an application.
- **PHP**. PHP Injection. Attack type in which a malicious PHP script is uploaded to a website. This often takes place by using a poorly-constructed upload form.
- **TROJAN**. Attack type in which malicious code poses as a legitimate app, script, or link, and tricks users into downloading and executing the malware on their local device.
- **DDOS**. (Direct Denial of Service). Attack type designed to bring down (or at least severely disrupt) a website. Typically, this is done by overwhelming the site with tens of thousands of spurious requests.
- **IN**. (Inbound Anomaly). Specifies the anomaly score of an inbound request. In anomaly scoring, requests aren't judged by a single rule; instead, multiple rules – and the past historical accuracy of those rules – determine whether or not a request is malicious.
- **OUT**. (Outbound Anomaly). Specifies the anomaly score of an outbound request.

Does any of this really matter to you? If you're using automated attack groups, yes, it really *does* matter. That's because, with automated attack groups, you don't manage individual rules; instead, you manage categories (i.e., attack groups). For example, to deny requests that violate an SQL Injection rule you set the SQL attack group `action` to **deny**; in turn, any request that triggers any rule in the group is denied. You don't have to know which rules are in the group, you don't have to know if new rules have been added or obsolete rules have been deleted, you just have to know how to set the SQL attack group to **deny**. In the Terraform block shown above, we set the attack group actions for all our attack groups to **deny**, which means that we're going to deny any request that triggers any group.

Admittedly, the block for configuring the attack group settings is a bit more complicated than most of the other Terraform blocks in this article. Because of that, we'll take a few minutes to walk you through the code, step-by-step.

The block starts off by calling the [akamai_appsec_attack_group](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_attack_group) resource; that's pretty routine. After that first line, however, we encounter this:

```
for_each = toset(["SQL", "XSS", "CMD", "HTTP", "RFI", "PHP", "TROJAN", "DDOS", "IN", "OUT"])
```

What's going on here? Here we set up a **for_each** loop that enables us to configure each individual attack group, one at a time. The attack group IDs (SQL, XSS, CMD, etc.) are configured as an array, and the **toset** function is used to convert this array into a set of strings, the data type required for use with a **for_eac**h loop.

Speaking of which, here's the code executed each time the loop is called:

```
config_id           = akamai_appsec_configuration.create_config.config_id
security_policy_id  = akamai_appsec_security_policy.security_policy_create.security_policy_id
attack_group        = each.value
attack_group_action = "deny"
```

The first two lines simply specify the IDs of the security configuration and the security policy; those IDs are required when working with the akamai_appsec_attack_group resource. That brings us to this line:

```
attack_group = each.value
```

In this line we specify the ID of the attack group being configured. However, we don't hardcode the ID; instead we use the syntax **each.value**. When we use a **for_each** loop, the first time through the loop the **each.value** property represents the first value in the loop; for us, that's **SQL**. That means that, the first time through the loop, we configure the SQL attack group. When that's done, we loop around, and **each value** now represents **XSS**, the second value in the **for_each** loop. That means we now configure the XSS attack group. This continues until we've looped through all the values included in the **for_each** loop.

And what exactly are we configuring? We're simply setting the `attack_group_action` for each attack group to **deny**. In other words, any request that triggers any attack group is denied.



### <a id="step16"></a>Step 16: Enable and configure the penalty box

```
resource "akamai_appsec_penalty_box" "penalty_box" {
 config_id              = akamai_appsec_configuration.create_config.config_id
 security_policy_id     = akamai_appsec_security_policy.security_policy_create.security_policy_id
 penalty_box_protection = true
 penalty_box_action     = "deny"
}
```

In ice hockey (and in a few other sports), players who commit more-egregious fouls are removed from the game and sent to the “penalty box.” Players in the penalty box, at least for purposes of the game, don't exist: they can't participate until they've served their time; in addition, and because they can't participate, they're pretty much ignored. Eventually players in the penalty box are allowed back into the game, but if they commit another foul they're returned to the penalty box and the cycle starts over again.

The Akamai penalty box serves a similar function: if a request triggers an attack group, the offending client is sent to the penalty box for 10 minutes. That means that all requests from that client are ignored during that 10-minute period. When time is up, the client can resume making requests, but another violation will send the client back to the penalty box for 10 more minutes.

> **Note**. OK, yes: there's a bit more nuance to the penalty box than what we've described here, but this is explanation enough for our purposes.

To employ the penalty box (available only if you're using automated attack groups), use the [akamai_appsec_penalty_box](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_penalty_box) resource and the following two arguments:

- `penalty_box_protection`. Set to **true** to enable the penalty box, or set to **false** to disable the penalty box.
- `penalty_box_action`. Set to **deny** to deny all requests from clients in the penalty box, or set to **alert** to issue an alert any time a client in the penalty box makes a request.

Note that the 10-minute timeout period is not configurable.

### <a id="step17"></a>Step 17: Activate the security configuration

```
resource "akamai_appsec_activations" "new_activation" {
  config_id           = akamai_appsec_configuration.create_config.config_id
  network             = "STAGING"
  notes               = "Activates the Documentation AAG Test Configuration on the staging network."
  activate            = true
  notification_emails = ["gstemp@akamai.com","karim.nafir@mail.com"]
}
```

When you create a security configuration, that configuration is automatically set to **inactive**; that means that the security configuration isn't actually analyzing and taking action on requests sent to your website. For that to happen, you employ the [akamai_appsec_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_activations) resource to activate the configuration. We won't go into the hows and whys of activating a security configuration in this documentation; for those details, see the **Activating a Security Configuration** article instead. Here we'll simply note that, when activating a configuration, you must specify the network where the configuration will be active. The akamai_appsec_activations resource gives you two choices when picking a network:

- **STAGING**. Typically, you start by activating a configuration on the staging network (as in the sample block shown above). The staging network consists of a small number of Akamai edge servers and provides a sandbox environment for testing and fine-tuning your configuration. Note that a configuration on the staging network idoesn't with your actual website and your actual website requests. Again, this network is for testing and fine-tuning, not for protecting your site.
- **PRODUCTION**. After you're satisfied with the performance of your security configuration, you can activate that configuration on the production network. Once activated there the configuration works with your actual website and your actual website requests.

Without going into too much detail, we activate our security configuration by using the akamai_appsec_activations resource and the following arguments:

| **Argument**          | **Description**                                              |
| --------------------- | ------------------------------------------------------------ |
| `config_id`           | Unique identifier of the configuration being activated.     |
| `network`             | Specifies which network the security configuration is being activated on. Allowed values are:     <br />* **staging**  <br />* **production** |
| `notes`               | Arbitrary notes regarding the network and its activation status. |
| `activate`            | If **true**, the specified network is activated; if **false**, the specified network is deactivated. Note that this property is optional: if omitted, `activate` is set to **true** and the network is activated. |
| `notification_emails` | JSON array of email addresses representing the people who receive a notification email when the activation process completes. |
