---
layout: "akamai"
page_title: "Creating Security Policies by Using Terraform"
description: |-
  Creating Security Policies by Using Terraform
---


# Creating Security Policies by Using Terraform

At heart, websites are all about handling user requests: this user would like to download a file; that user would like to stream a video; a third user would like to visit one of your web pages. The vast majority of these requests are legitimate, and harmless; in fact, they're the very reason you published your website in the first place. However, other requests (either maliciously or inadvertently) might not be so harmless. To safely and securely manage your website, you need to be able to identify suspect requests;and to quickly and efficiently deal with these requests.

At Akamai, security policies play a key role in identifying and handling website requests. If a request is flagged by a match target (that is, if the request matches criteria you have specified in advance) the security policy associated with that match target can step in and provide a more detailed analysis on the request, applying protections such as rate limiting and reputation controls to help verify the legitimacy and the safety of the request. Requests that pass these tests are allowed through; depending on how you have configured your policies, requests that don't pass these tests can be rejected.

That, in a nutshell, is why you need security policies.

And by using Terraform, there are at least three different ways you can create a security policy. You can:

- Create a new security policy that uses the default policy settings (these settings will be discussed later in this documentation).

- Create a “cloned” security policy: a new policy that inherits the settings and setting values assigned to an existing policy. For example, cloned security Policy B will be, in effect, an exact duplicate of existing security Policy A.

- Create a new security policy and, in the same Terraform operation, configure custom values for the security policies settings.


#### A Note About Using Multiple Security Policies

A single security configuration can have multiple security policies. And that's good: for example, you might have one set of APIs that require a different set of protections than your other APIs. In that case, you might need two security policies: one for the “special” set of APIs and the other for the remaining APIs. What you probably don't need, however, is one security policy for each individual API. Having multiple security policies provides you with flexibility and with the opportunity to fine-tune your protections. At the same time, however, each security policy you add to a security configuration increases the time it take to analyze and process each request. You'll need to find a balance between having customized protections and between having an efficient and responsive website. As a general rule, the fewer security policies you need the better.

## Security Policy Default Settings

Unless you specify otherwise, any new security policy you create will be assigned the default policy settings. The current values for these settings can be returned by using the [akamai_appsec_security_policy_protections](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_security_policy_protections) data source and a Terraform configuration similar to this:

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
data "akamai_appsec_security_policy_protections" "protections" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "protections_response" {
  value = data.akamai_appsec_security_policy_protections.protections.output_test
}
```

As you can see, this is a pretty straightforward configuration file: we call the Akamai provider and point that provider to our authentication credentials (stored in the **.edgerc** file), We connect to the Documentation configuration, then use this block to return protections information for the security policy **gms1_134637**:

```
data "akamai_appsec_security_policy_protections" "protections" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}
```

After the protections data is returned, we then use the last block in the file to output that information to the screen:

```
output "protections_response" {
  value = data.akamai_appsec_security_policy_protections.protections.output_test
}
```

When all is said and done you'll get back response similar to this:

```
+------------------------------------------------------------------------------------------------------------------------------------------+
| wafProtectionDS                                                                                                                          |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| APICONSTRAINTS | APPLICATIONLAYERCONTROLS | BOTMANCONTROLS | NETWORKLAYERCONTROLS | RATECONTROLS | REPUTATIONCONTROLS | SLOWPOSTCONTROLS |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| true           | true                     | false          | true                 | true         | true               | true             |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
```

These properties are described in the following table:

| API Property             | Description                                                  | Default Value |
| ------------------------ | ------------------------------------------------------------ | ------------- |
| APICONSTRAINTS           | Places limits on both the number and the size of API requests sent by a given user. |               |
| APPLICATIONLAYERCONTROLS | Uses the Web Application Firewall (WAF) to help minimize the effects of cross-site scripting, SQL injection, file inclusion, and other attacks. | true          |
| BOTMANCONTROLS           | Places limits on the number of valid, and invalid, form submissions associated with a single user. | false         |
| NETWORKLAYERCONTROLS     | Block (or allows) requests based in a client's IP address or geographic location. | true          |
| RATECONTROLS             | Provides a way to monitor, and to control, the rate of requests received by your site. | true          |
| REPUTATIONCONTROLS       | Helps identify potentially-malicious clients based on past behaviors associated with the client's IP address. | false         |
| SLOWPOSTCONTROLS         | Helps guard against Denial of Service attacks caused by extremely slow request rates. | true          |

When you enable one of these features, that feature is automatically assigned a set of default values,. For example, if you enable slow POST controls then, by default, the feature is configured with the following settings:

| Action | SLOW_RATE_THRESHOLD<br /> RATE | SLOW_RATE_THRESHOLD<br /> PERIOD | DURATION_THRESHOLD<br />TIMEOUT |
| ------ | ------------------------------ | -------------------------------- | ------------------------------- |
| alert  | 10                             | 60                               | Null                            |


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

## Creating a Security Policy that Uses the Default Settings

The quickest and easiest way to create a security policy is to create a policy that uses the default settings. As a reminder, those defaults are shown below:

| Setting                  | Default Value |
| ------------------------ | ------------- |
| APICONSTRAINTS           | true          |
| APPLICATIONLAYERCONTROLS | true          |
| BOTMANCONTROLS           | false         |
| NETWORKLAYERCONTROLS     | true          |
| RATECONTROLS             | true          |
| REPUTATIONCONTROLS       | false         |
| SLOWPOSTCONTROLS         | true          |

And here's a sample Terraform collection that creates a security policy that uses those default settings:

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
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  default_settings       = true
  security_policy_name   = "New Default Policy"
  security_policy_prefix = "gms2"
}

output "security_policy_create" {
  value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

As you can see, this is another straightforward Terraform configuration: it simply creates a new security policy (with the name **New Default Policy** and the prefix **gms2**) and associates that policy with the **Documentation** security configuration. Note that you don't include an ID for the security policy when creating your configuration file: IDs are assigned by the system when the policy is created. That ID will comprise the security policy prefix, an underscore (_) and a numeric value assigned by Akamai. For example:

```
gms2_135566
```

Note, too that our Terraform configuration also includes this argument:

```
default_settings = true
```

This simply tells Terraform that we want the new security policy to use the default security policy settings. The default_settings argument optional: if it's not included the policy will automatically be assigned the default settings. And if you set `default_settings`to **false**? In that case, all the policy settings are also set to **false**:

```
+------------------------------------------------------------------------------------------------------------------------------------------+
| wafProtectionDS                                                                                                                          |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| APICONSTRAINTS | APPLICATIONLAYERCONTROLS | BOTMANCONTROLS | NETWORKLAYERCONTROLS | RATECONTROLS | REPUTATIONCONTROLS | SLOWPOSTCONTROLS |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
| false          | false                    | false          | false                | false        | false               | false          |
+----------------+--------------------------+----------------+----------------------+--------------+--------------------+------------------+
```

After your Terraform configuration is ready, you should run the `terraform plan` command from the command prompt. Running the plan command does two things for you. First, it does some syntax checking, and alerts you to many (although not necessarily all) configuration errors:

```
│ Error: Missing required argument
│
│   on akamai.tf line 17, in resource "akamai_appsec_security_policy" "security_policy_rename":
│   17: resource "akamai_appsec_security_policy" "security_policy_rename" {
│
│ The argument "security_policy_prefix" is required, but no definition was found.
```

Second, it tells you exactly what will happen if you create the security policy:

```
Terraform will perform the following actions:

akamai_appsec_security_policy.security_policy_create will be created

resource "akamai_appsec_security_policy" "security_policy_create" {
  config_id              = 58843
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
akamai_appsec_security_policy.security_policy_create: Creation complete after 9s [id=58843:gms2_135566]

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
akamai_appsec_security_policy.security_policy_create: Creation complete after 9s [id=58843:gms2_135566]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

> **Hint**. It's the **gms2_135566** portion of **id=58843:gms2_135566**. That value happens to be the ID of the security configuration (**58843**) followed by a colon followed by the security policy ID.

## “Cloning” a Security Policy

To clone something means to make an exact replica of that something, and that's exactly what happens when you clone a security policy: you create a new policy (Policy B) that has the exact same settings and setting values as an existing policy (Policy B). Other than their identifiers (such as the policy name and the policy ID), the two policies will be indistinguishable from one another.

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

The preceding configuration creates a new security policy (a policy with the name **Cloned Policy** and the prefix **gms3**) as part of the **Documentation** security configuration. Unlike our previous example, this new policy isn't based on the default settings; instead, the policy is, for all intents and purposes, a duplicate of the existing security policy **gms1_134637**. That simply means that the policy will be configured using the exact same settings and setting values assigned to gms1_134637. In other words:

| API Property             | gms_134637 | Cloned Policy |
| ------------------------ | ---------- | ------------- |
| APICONSTRAINTS           | false      | false         |
| APPLICATIONLAYERCONTROLS | false      | false         |
| BOTMANCONTROLS           | false      | false         |
| NETWORKLAYERCONTROLS     | true       | true          |
| RATECONTROLS             | true       | true          |
| REPUTATIONCONTROLS       | false      | false         |
| SLOWPOSTCONTROLS         | true       | true          |

In order to clone a security policy, your Terraform configuration must include the following two arguments:

```
default_settings = false
create_from_security_policy_id = "gms1_134637"
```

The first argument (**default_settings = false**) simply tells Terraform *not* to apply the default settings to the new policy. Instead, we want our new policy to have the same settings and settings values as the ones assigned to the existing security policy gms1_134637. As you might have guessed, that's what the second argument does: it specifies the ID of the security policy whose settings we want copied to the new policy.

FThe rest is easy. Just like we did before, run `terraform plan` to verify your syntax, and then run `terraform apply` to create the new security policy. If everything goes according to plan, the policy will be created, and you'll see output similar to this:

```
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
Outputs
security_policy_create = "gms3_135568"
```

## Adding Custom Setting Values When Creating a Security Policy

As alluded to previously, another option available when creating a security policy is to create a policy that includes custom setting values; for example, you can create a policy that enables slow post protection, but that doesn't use the default setting values:

| Slow Post Setting          | Default Value | New Policy Value |
| -------------------------- | ------------- | ---------------- |
| action                     | alert         | abort            |
| SLOW_RATE_THRESHOLD RATE   | 10            | 15               |
| SLOW_RATE_THRESHOLD PERIOD | 60            | 30               |
| DURATION_THRESHOLD TIMEOUT | null          | 20               |

To do this, we'll use a Terraform configuration similar to this:

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

In this case, however, the configuration doesn't end once the new policy has been created. Instead, it continues on to assign custom values to the slow POST control settings:

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

> b. In this case, there was no need to enable slow POST controls. That's because the policy was creating using the default settings and, by default, slow POST control is enabled in all new security policies.
>

In the preceding block, the only “tricky” line is this one:

```
security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
```

In this line we specify the ID of the security policy we want to update. Of course, this is a brand-new security policy which didn't even exist few seconds ago. So how are we supposed to know the ID of the new policy? By using this value:

```
akamai_appsec_security_policy.security_policy_create.security_policy_id
```

This is the value we echo back to the screen immediately after the security policy has been created:

```
output "security_policy_create" {
  value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

Needless to say, that's also the ID of the new policy. Note the final line in this snippet from the Terraform response when creating a new policy:

```
akamai_appsec_security_policy.security_policy_create: Creating...
akamai_appsec_security_policy.security_policy_create: Creation complete after 8s [id=58843:gms4_135620]
```

In other words, what we do here is:

1. Create a new security policy, one that uses the default settings.
2. Grab the ID from that new policy and then use that ID to modify the slow post control settings.

That's all there is to it. Assuming everything works, we should see the following output as our Terraform configuration completes:

```
akamai_appsec_security_policy.security_policy_create: Creating...
akamai_appsec_security_policy.security_policy_create: Creation complete after 8s [id=58843:gms4_135620]
akamai_appsec_slow_post.slow_post: Creating...
akamai_appsec_slowpost_protection.protection: Creating...
akamai_appsec_slowpost_protection.protection: Creation complete after 4s [id=58843:gms4_135620]
akamai_appsec_slow_post.slow_post: Creation complete after 4s [id=58843:gms4_135620]

Apply complete! Resources: 3 added, 0 changed, 0 destroyed.

Outputs:
security_policy_create = "gms4_135620"
```
