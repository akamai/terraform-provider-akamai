---
layout: "akamai"
page_title: "Creating a Security Configuration"
description: |-
  Creating a Security Configuration
---


# Creating a Security Configuration

In the Akamai world, a security configuration primarily functions as a container for other application security objects: security policies, rate policies, reputation profiles, Kona Site Defender rules, etc. Depending on your needs, you might have multiple security configurations: many organizations have different security configurations for their different business units, their different domains, their different geographic units, and so on. If you do need to create multiple configurations then you're in luck: Terraform provides three different ways for you to create new security configurations. You can:

- Create a “blank” security configuration
- Create a security configuration that uses the default values
- Create a new security configuration by “cloning” an existing configuration

## Creating a Blank Security Configuration

In Control Center, the easiest way to create a security configuration is to select the option **Manually (Create a blank Web Security Configuration)**. This approach is easy because you don't have do much beyond specifying:

- The contract and group ID associated with the new configuration.
- A name and (optionally) description of the new configuration.
- At least one “selectable hostname.” You can't create a security configuration without including at least one host to be protected by that configuration.

That's how you do it in Control Center, and that's also how you do things in Terraform. If you go this route, just keep in mind that your cinfiguration really *will* be empty: the new configuration will contain at least one host, but that's it. To make your configuration useful, you'll need to add security policies and match targets and rate policies and …. Be that as it may, if you want to start with the basics and then begin adding on your new configuration from there, creating a blank (empty) security configuration doesn't require too much effort on your part. In fact, the following Terraform configuration is all that's needed to create a blank configuration named **Empty Security Configuration**:

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

This, as you probably figured out for yourself, is where we use the akamai_appsec_configuration resource to create the new configuration. Creating a blank configuration requires us to supply the following five arguments:

| Argument    | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| name        | Unique name to be assigned to the new configuration.         |
| description | Brief description of the configuration and its intended purpose. |
| contract_id | Akamai contract ID associated with the new configuration. You can use the [akamai_appsec_contracts_groups](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_contracts_groups) data source to return information about  the contracts and groups available to you. |
| group_id    | Akamai group ID associated with the new configuration.       |
| host_names  | Names of the selectable hosts to be protected by the configuration. Note that names must be passed as an array; that's what the square brackets surround "documentation.akamai.com" are for. If you want to add multiple hostnames to the configuration, just separate the individual names by using commas. For example:<br /><br />`host_names = ["documentation.akamai.com", "training.akamai.com", "events.akamai.com"]` |

Finally, we use this block to echo back the ID of the newly-created security configuration:

```
output "create_config_id" {
  value = akamai_appsec_configuration.create_config.config_id
}
```

When we call this configuration, we should get output similar to the following:

```
akamai_appsec_configuration.create_config: Creating...
akamai_appsec_configuration.create_config: Creation complete after 5s [id=76967]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

create_config_id = 71219
```

## Creating a Security Configuration That Uses the Recommended Presets

As noted previously, when you create an empty security configuration you get just that: a security configuration that doesn't include much of anything (including a security policy). That means that you'll have to do some work before you can actually start to *use* that security configuration.

In some cases, however, you might want to get a security configuration up and running as quickly as possible, and then fine-tune and adjust it later. One easy way to do that is to create a security configuration that uses the recommended (default) settings, When you take this route you end up with a security configuration that:

- Has a single security policy that uses the default values.
- Has the IP/Geo Firewall enabled (although that's only a convenience, because you won't have any network lists to be allowed or blocked by the firewall).
- Has three rate policies: **Origin Server**; **POST Page Requests**; **Page View Requests**. In all three policies both the IPv4 and the IPv6 actions are set to **alert**.
- Has slow post protection enabled using the default values. This includes setting the slow post action to **alert**.
- Has no custom rules.
- Enables you to choose between Automated Attack Groups (**AAG**, in which firewall rules are automatically updated and maintained for you), and Kona Rule Set (**KRS**, in which you have the ability to configure each rule's action, conditions, and exceptions).
- Has API request constraints enabled and the action set to **alert**. However, no API matches are defined.
- Has single match target: All Hostnames (which matches the path **/***).

In this documentation, we won't create a configuration that includes all the default items. Instead, we'll create a simpler configuration that: 1) includes a security policy that uses the default security policy settings; and, 2) sets the Web Application Firewall (WAF) mode to **KRS** (Kone Rule Set). Although not a fully-fleshed out configuration, this should give you an idea how to:

- Create a security configuration.
- Use the ID of the new configuration to create a security policy.
- Use the ID of the security policy to modify the WAF mode.

If you can master the preceding tasks then creating a configuration that uses all the default settings shouldn't be too hard. For example, here's our sample Terraform configuration:

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

After declaring the Akamai provider and providing our authentication credentials, we use the akamai_appsec_configuration resource and this block to create the security configuration:

```
resource "akamai_appsec_configuration" "create_config" {
  name        = "Documentation Security Configuration"
  description = "This security configuration is used by the documentation team."
  contract_id = "1-3CV382"
  group_id    = 47346
  host_names  = ["ldap.host1.akamai.com.edgesuite-staging.net"]
}
```

If this looks familiar, well, that shouldn't come as a surprise: it's the exact same Terraform configuration we used when creating an empty security configuration. The only difference is that, when we created our empty configuration, we stopped right here. This time, we're going to keep going.

After the security configuration is complete, and after we know the ID of that configuration, we can then create our security policy. But how are we supposed to know the ID of a brand-new security configuration?

As it turns out, determining the ID of our new configuration is surprisingly easy. To begin with, after a new configuration has been created, the ID of that configuration is available to us in an attribute named **config_id**; we know that because that information is included in the akamai_appsec_configuration resource documentation.

On top of that, the first line in our resource block looks like this:

```
resource "akamai_appsec_configuration" "create_config" {
```

As you know, **akamai_appsec_configuration** is the name of the resource used to create a configuration. Meanwhile, **create_config** is an attribute that represents the newly-created configuration. That means we have three pieces of information to work with:

| Item                        | Description                                                  |
| --------------------------- | ------------------------------------------------------------ |
| akamai_appsec_configuration | The Terraform resource.                                      |
| create_config               | The attribute that references the new configuration.         |
| config_id                   | The attribute that contains the ID of the new configuration. |


To reference the ID of our new configuration, we just need to use “dot notation” to string these three items together:

```
akamai_appsec_configuration.create_config.config_id
```

If you look at the Terraform block that creates a security policy, you'll see that we use the preceding as the value of the `config_id` argument:

```
resource "akamai_appsec_security_policy" "security_policy_create" {
  config_id              = akamai_appsec_configuration.create_config.config_id
  default_settings       = true
  security_policy_name   = “"Documentation Security Policy"
  security_policy_prefix = "doc1"
}
```

> **Note**. We won't go into any of the other details involved in creating a security policy. You can find that information in the article Creating a Security Policy.

At this point we have a new configuration and a new security policy; all that's left is to specify the WAF mode. To do that, we need two things: the configuration ID and the security policy ID. We already have the configuration ID: we just used it to create the security policy. But what about the security policy ID?

Have no fear: any time you create a new security policy an attribute named **security_policy_id** is made available to you. And, if you look at the first line of our security policy block, you'll see both a resource name (**akamai_appsec_security_policy**) and an attribute that represents the new policy (**security_policy_create**). How do we know the ID of our new security policy? Once again, we use dot notation to string together our individual elements:

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

All that's left now is to set the `mode` to **KRS**.

And that's it (at least for our immediate purposes). When we run the `terraform apply` command we'll get back a response similar to the following:

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

If you take a closer look at the response, you'll see the ID of the new configuration (**76984**) and the ID of the new security policy (**doc1_137405**). Success!

## Cloning a Security Configuration

Another way to create a new – and functional – security configuration is to make a “clone” of an existing configuration: in effect, you take security configuration A and create an exact (or at least a nearly-exact) replica of that configuration (e.g., Configuration B). For example, among the items in Configuration A that will be replicated in Configuration B are:

- Security policies
- Rate policies
- Custom rules
- Custom denies
- SIEM settings
- Advanced logging and prefetch settings
- Slow post settings
- Match targets

Etc., etc. Note that, in most cases, items are copied exactly the way they appear in Configuration A. For example, suppose Configuration A contains the following security policies:

- Security Policy 1
- Security Policy 2
- Security Policy 3
- Security Policy 4

After the cloning operation is complete, Configuration B will contain the exact same set of policies, with the exact same names and setting values:

- Security Policy 1
- Security Policy 2
- Security Policy 3
- Security Policy 4

In other cases the Akamai provider will take into account the fact that the new configuration is a, well, new configuration. For example, in Configuration A the Web Security Configuration ID setting is shown as 58843.

When Configuration A is cloned, however, the Web Security Configuration ID setting will reflect the ID of the new configuration (Configuration B).

Your Terraform configuration for cloning a security configuration will look similar to this:

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

resource "akamai_appsec_configuration" "clone_config" {
  name                  = "Cloned Security Configuration"
  description           = "This security configuration is based on configuration ID 58843."
  create_from_config_id = 58843
  create_from_version   = 9
  contract_id           = "1-3UW382"
  group_id              = 13139
  host_names            = ["documentation.akamai.com"]
}

output "clone_config_id" {
  value = akamai_appsec_configuration.clone_config.config_id
}
```

As usual, there's nothing especially complicated about the configuration: this one starts by declaring the Akamai Terraform provider, and by presenting our authentication credentials. From there, we move directly into the block that uses the akamai_appsec_configuration resource to clone configuration 58843:

```
resource "akamai_appsec_configuration" "clone_config" {
  name                  = "Cloned Security Configuration"
  description           = "This security configuration is based on configuration ID 58843."
  create_from_config_id = 58843
  create_from_version   = 9
  contract_id           = "1-3UW382"
  group_id              = 13139
  host_names            = ["documentation.akamai.com"]
}
```

The akamai_appsec_configuration block for cloning a configuration is remarkably similar to the other blocks we've looked at (such as the block for creating a blank configuration). The only difference is that, when cloning a configuration, we need to include two additional arguments:

```
create_from_config_id = 58843
create_from_version   = 9
```

The `create_from_config_id` argument specifies the ID of the security configuration you want to replicate (in our original analogy, Configuration A). Similarly, the `create_from_version` enables you to select a specific version of the configuration to be cloned. Note that this argument is optional: if it's not included, Terraform will replicate the initial (v1) version of the configuration.

The final block of code is optional: it simply echoes back the ID of the new security configuration:

```
output "clone_config_id" {
  value = akamai_appsec_configuration.clone_config.config_id
}
```

When you run the terraform apply command, you should see output similar to this:

```
akamai_appsec_configuration.clone_config: Creating...
akamai_appsec_configuration.clone_config: Still creating... [10s elapsed]
akamai_appsec_configuration.clone_config: Creation complete after 11s [id=76982]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

clone_config_id = 76982
```

In this example, the new configuration has been assigned the ID **76982**.

## Activating a Security Configuration

By default, security configurations aren't activated when they're created; that simply means that those configurations aren't actually analyzing and responding to requests. To actually  use of a security configuration, that configuration needs to be activated. Typically, that's a two-step process: first the configuration is activated on the staging network and then, after testing and fine-tuning, the configuration is activated on the production network. At that point, the configuration is fully deployed, and is analyzing and responding to requests.

When the time comes, you can use Terraform to activate a configuration to either the staging network or the production network. A Terraform configuration for doing this will look similar to the following:

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

resource "akamai_appsec_activations" "activation" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  network             = "STAGING"
  notes               = "This is a test configuration used by the documentation team."
  notification_emails = ["gstemp@akamai.com"]
}
```

For the most part, this is a typical terraform configuration: we declare the Akamai provider, provide our authentication credentials, and connect to the Documentation configuration. We then use the [akamai_appsec_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_activations) resource and the following block to activate that configuration:

```
resource "akamai_appsec_activations" "activation" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  network             = "STAGING"
  notes               = "This is a test configuration used by the documentation team."
  notification_emails = ["gstemp@akamai.com"]
}
```

Inside this block we need to include the following arguments and argument values:

| Argument            | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| config_id           | Unique identifier of the configuration being activated.      |
| network             | Name of the network the configuration is being activated for. Allowed values are:<br /><br />* STAGING<br />* PRODUCTION |
| notes               | Information about the configuration and its activation.      |
| notification_emails | JSON array of email addresses for people who should be notified when activation is complete. To send notification emails to multiple people, separate the individual email addresses by using commas:<br /><br />notification_emails = ["gstemp@akamai.com", "karim.nafir@mail.com"] |


From here we can run `terraform plan` to verify our syntax, then run `terraform apply` to activate the security configuration.
