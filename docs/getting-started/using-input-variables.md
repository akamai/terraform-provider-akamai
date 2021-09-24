---
layout: "akamai"
page_title: "Using Input Variables in a Terraform Configuration"
description: |-
  Using Input Variables in a Terraform Configuration
---


# Using Input Variables in a Terraform Configuration

In the Application Security Terraform user guides, we hard-code values for such things as security configuration ID or security policy ID. That's why you see things like **58843** and **gms1_134637** in our code samples:

```
resource  "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
  config_id = 58843
  security_policy_id = "gms1_134637"
  rate_policy_id = 122149
  ipv4_action = "alert"
  ipv6_action = "alert"
}
```

As a general rule, hard-coding values makes it easier for people to understand what's going on, and what information they'll need to supply in order to modify a code sample. In other words, we hard-code values for educational reasons.

However, that's not necessarily the way that people who write Terraform configurations for a living do things. For example, if you look at other Terraform samples you might see coding that looks more like this:

```
resource  "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
  config_id = var.configuration_id
  security_policy_id = var.security_policy_id
  rate_policy_id = var.rate_policy_id
  ipv4_action = var.set_alert
  ipv6_action = var.set_alert
}
```

And that's exactly what we thought when we first encountered a Terraform code sample: what in the world is a **var.configuration_id** or a **var.security_policy_id**?

As it turns out, all those constructions that begin with **var.** are Terraform variables (or, to be a little more precise, Terraform input variables). If you know anything at all about computer programming, or if you know anything at all about basic algebra, then you're no doubt familiar with variables. For example, take the following equation:

```
5 + X = 8
```

In the preceding equation, **X** isn't a number: there is no number X. Instead, **X** is a variable that *represents* a number. In this case, that's the number 3; in other cases, X represents a different number:

```
5 + X = 17
```

They're called variables because – well, because their values can vary.

Variables in Terraform serve a similar purpose; for example, the variable **var. configuration_id** represents the ID of a specific security configuration. When we write a Terraform configuration we can hard-code the configuration ID:

```
config_id = 58843
```

Or, we can use a variable, in which case Terraform will use the value assigned to **var.configuration_id** as the configuration ID:

```
config_id = var.configuration_id
```

Either way, we'll end up working with security configuration 58843.

## So Does This Mean I Should Use Variables Rather Than Hard-Coded Values?

That's up to you; like we said, you get the same results either way (because a hard-coded **58843** is exactly the same as a **58843** stored in a variable). Nevertheless, variables do offer a few advantages that hard-coded values don't. For one thing, configurations that use variables are more portable than configurations that rely on hard-coded values. Without making changes, the following configuration can only be used by someone working with configuration 58843:

```
config_id = 58843
```

By comparison, this configuration below can be used -- unchanged -- by anyone who has predefined the variable **var.configuration_id** (we'll have more on that in a minute):

```
config_id = var.configuration_id
```

Variables also allow you to write code without having to memorize (or to constantly look up) certain values. Do you have any idea what your contract ID and Akamai group ID are? Don't feel bad: we don't know ours, either. But that's OK; after all, any time we need to specify our contract and group IDs we just use variables :

```
contract_id = var.contract_id
group_id = var.group.Id
```

Variables are also useful in validating and troubleshooting your Terraform configurations. But that's another thing we'll discuss in a minute. First we need to talk about how you can define and make use of Terraform variables.

## Defining a Terraform Variables

Although there different ways to define a Terraform variable, we'll focus here on a very basic approach: creating a single .tf file to house your variable definitions. (Yes, there are other ways to work with variables, but we're going to keep things easy for the moment.) Note that you can give your variable file any name you want, as long as that file has a **.tf** file extension and as long as that file is in the same folder as your Terraform executable.

> **Note**. In case you're wondering, we gave our variables file the not-so-clever name **variables.t**f.

In the variables file, we then have separate definition blocks for each variable. For example, here's the definition block for the variable configuration_id:

```
variable "configuration_id" {
  default        = 58843
  description    = "Unique identifier of a security configuration. Unlike security configuration names, the configuration ID can't be changed."
}
```

Before we go any further, yes, the variable name is **configuration_id**; however, when writing a Terraform configuration file we reference that variable as **var.configuration_id**. That's just the variable named preceded by **var.**, which indicates that we're using a variable.

You might have noticed that our variable definition has two arguments:

- `default`. This is the value assigned to the variable by default. For example, take the following snippet from a Terraform configuration file:

  `config_id = var.configuration_id`

  Because the variable **var.configuration_id** has the default value **58843**, this snippet connects us to security configuration 58843.

- `description`. A brief description of the variable and what it represents. The description is optional, but it *can* be useful, especially if you're defining a large number of variables or other people will be sharing your variables file.

That's all we need, at least for the moment. Let's now take a look at a very simple Terraform configuration that uses our newly-defined variable. This particular configuration returns versioning information for the specified security configuration:

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

data "akamai_appsec_configuration_version" "versions" {
  config_id = var.configuration_id
}

output "versions_output_text" {
  value = data.akamai_appsec_configuration_version.versions.output_text
}
```

If we run the `terraform plan` command, we'll get back versioning information for configuration 58843:

```
+-----------------------------------------------------+
| configurationVersion                                |
+----------------+----------------+-------------------+
| VERSION NUMBER | STAGING STATUS | PRODUCTION STATUS |
+----------------+----------------+-------------------+
| 1              | Inactive       | Inactive          |
| 2              | Inactive       | Inactive          |
| 3              | Inactive       | Inactive          |
| 4              | Inactive       | Inactive          |
| 5              | Inactive       | Inactive          |
| 6              | Inactive       | Inactive          |
| 7              | Inactive       | Inactive          |
| 8              | Inactive       | Inactive          |
| 9              | Inactive       | Inactive          |
+----------------+----------------+-------------------+
```

Which is exactly what we wanted to get back.

And what if we wanted to get back versioning information for configuration 76478 instead? That's fine: one way to do that would be to change the default value for **configuration_id** in the variables file. We'll show you another way momentarily.

## Specifying the Datatype for a Variable

It's not unusual to mistake a configuration ID (which has a numeric value like **55843**) and a configuration name (which has a string value like **Documentation**). If you enter the configuration name where you  should have entered the configuration ID (i.e., you used a string value instead of a numeric value) your configuration is bound to fail. And, depending on the nature of the failure, you might find this issue difficult to troubleshoot.

One way that variables can help minimize this kind of mistake is by enabling you to specify the datatype for a variable. Variables can use one of the following data types:

- **string**
- **number**
- **bool** (Boolean)

For example, a revised version of our variable definition might look similar to this:

```
variable "configuration_id" {
  default        = 58843
  description    = “Unique identifier of a security configuration. Unlike security configuration names, the configuration ID can't be changed.”
  type           = number
}
```

How does this help? Well, suppose you set the value of `configuration_id` to **Documentation** when you call `terraform plan` or `terraform apply`. According to the revised variable definition, `configuration_id` has to be a number; needless to say, the value **Documentation** is *not* a number. As a result, our command fails with the following error:

```
╷
│ Error: Invalid default value for variable
│
│   on variables.tf line 287, in variable "configuration_id":
│  287:   default        = "Documentation"
│
│ This default value is not compatible with the variable's type constraint: a number is required.
╵
```

Yes, the command failed: it was bound to fail. But thanks to the `type` property, we know exactly why it failed.

## Adding a Validation Rule to a Variable

Specifying the datatype helps you  minimize mistakes, but specifying a datatype can only take you so far: after all, **1** is a number, and so is **83277563874583724582342293**. However, neither of those numbers represent a valid configuration ID. Ideally, we'd like to add some code that would provide users with a little more information as to what a valid configuration ID looks like.

As it turns out, that's what the `validation` property is for. For example, this updated variable definition not only sets the datatype, it also states that a configuration ID must have exactly 5 digits:

```
variable "configuration_id" {
  default        = 58843
  description    = "Unique identifier of a security configuration. Unlike security configuration names, the configuration ID can't be changed."
  type = number

  validation {
    validation {
    condition     = length(tostring(var.configuration_id)) == 5
    error_message = "The configuration ID must contain exactly 5 digits."  }
}
```

We won't provide a full explanation of variable validation in this document: you can read more about variables and validation in the article [Input Variables](https://www.terraform.io/docs/language/values/variables.html). For now, we'll simply dissect our sample validation:

```
  validation {
    validation {
    condition     = length(tostring(var.configuration_id)) == 5
    error_message = "The configuration ID must contain exactly 5 digits."  }
```

As you can see, the validation section contains two components:

- `condition`. The criterion that must be met for the variable value to pass validation. In this example, we're saying that the length of the value (**var.configuration_id**) must be equal to (note the back-to-back equal signs) **5**; that simply means that the value must contain *exactly* 5 characters, no more and no less. Oh, and because the length function works only on string values, we also use **tostring** to convert the numeric ID to a string.

  But don't get too hung up on the specifics: just remember that the condition is the criterion that must be met for the variable value to pass validation.

- `error_message`. Not too surprisingly, this is the message displayed to the user if the variable fails validation.

To use the validation all we have to do is call `terraform plan` or `terraform apply`. If we've supplied a valid configuration ID (like **58843**), we'll get back versioning information:

```
+-----------------------------------------------------+
| configurationVersion                                |
+----------------+----------------+-------------------+
| VERSION NUMBER | STAGING STATUS | PRODUCTION STATUS |
+----------------+----------------+-------------------+
| 1              | Inactive       | Inactive          |
| 2              | Inactive       | Inactive          |
| 3              | Inactive       | Inactive          |
| 4              | Inactive       | Inactive          |
| 5              | Inactive       | Inactive          |
| 6              | Inactive       | Inactive          |
| 7              | Inactive       | Inactive          |
| 8              | Inactive       | Inactive          |
| 9              | Inactive       | Inactive          |
+----------------+----------------+-------------------+
```

But what happens if we supply an invalid configuration ID, like the 6-digit value **123456**? This is what happens:

```
│ Error: Invalid value for variable
│
│   on variables.tf line 286:
│  286: variable "configuration_id" {
│
│ The configuration ID must contain exactly 5 digits.
│
│ This was checked by the validation rule at variables.tf:291,5-15.
```

## Calling Variables from the Command Line

Up to this point we've seen several advantages to using variables rather than hard-coded values. However, we're still faced with one problem. If we use hard-coded values, we have to change our Terraform configuration file any time we want to switch from security configuration 58843 to security configuration 76478. On the other hand, if we use variables we have to change the default value of the **configuration_id** variable if we want to switch security configurations. Either way, we have to change *something* any time we go from one security configuration to another.

Don't we?

Believe it or not, the answer to that question is this: no. That's because you don't have to use the default value assigned to a variable; in fact, the default value is used only when you don't supply a different value when calling `terraform plan` or `terraform apply`. What does that mean? Well, consider this sample call to `terraform plan`:

```
terraform plan -var='configuration_id=76478'
```

As you can see, this command includes the optional **-var** argument:

```
-var='configuration_id=76478'
```

As the names implies (sort of), the **-var** argument enables us to specify the name and value for a Terraform variable *at the time we run our configuration*. In this case, we're saying that we want to set the value of the **configuration_id** variable to 76478. Yes, in our variable definition we set the default value of **configuration_id** to 58843. Thanks to the **-var** argument, however, Terraform will override that default and use 76479 as the configuration ID. It's as though our Terraform configuration actually looked like this:

```
 data "akamai_appsec_configuration_version" "versions" {
  config_id = 76478
}
```

In response, we get back configuration versioning information for security configuration 76478:

```
+-----------------------------------------------------+
| configurationVersion                                |
+----------------+----------------+-------------------+
| VERSION NUMBER | STAGING STATUS | PRODUCTION STATUS |
+----------------+----------------+-------------------+
| 1              | Inactive       | Inactive          |
+----------------+----------------+-------------------+
```

When you enter a variable at the command prompt, the default value configured in the variable definition is ignored. However, everything else defined for the variable – such as any validation rules – are *not* ignored. For example, suppose we set configuration_id to a 6-digit value:

```
terraform plan -state=terraform.tfstate -var='configuration_id=123456'
```

Do that, and the validation rule will cause our command to fail:

```
│ Error: Invalid value for variable
│
│   on variables.tf line 286:
│  286: variable "configuration_id" {
│
│ The configuration ID must contain exactly 5 digits.
│
│ This was checked by the validation rule at variables.tf:291,5-15.
```

Pretty cool.
