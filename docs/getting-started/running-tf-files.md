---
layout: "akamai"
page_title: "Running Terraform Configuration (.tf) Files"
description: |-
  Running Terraform Configuration (.tf) Files
---


# Running Terraform Configuration (.tf) Files

To use Terraform, you create one or more Terraform configurations: text files (with a **.tf** file extension) written in the [HashiCorp Configuration Language](https://www.terraform.io/docs/language/index.html) (HCL), These text files (which can be written using any editor capable of saving files in plain-text format) look similar to this:

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
  name = var.security_configuration
}

data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

resource "akamai_appsec_eval_protect_host" "protect_host" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = ["documentation.akamai.com"]
}
```

> **Note**. For the moment, at least, don't worry about the content of this file or what the individual lines of code actually do. We'll talk about the HCL language and how you go about creating a Terraform configuration later in this documentation.

A Terraform configuration file contains the information required to do something to your Akamai infrastructure: create a security configuration, delete a rate policy, update your slow post protection settings. And this, of course, is what Terraform is all about: you add some configuration information to a text file and, in turn, you create, modify, or delete items in the Akamai infrastructure.

Of course, merely creating a Terraform configuration file doesn't mean that you've magically created a security policy or updated your Web Application Firewall settings. Instead, you also need to use Terraform to read and carry out the instructions included in that file. Terraform provides multiple ways to execute a Terraform configuration, including:

- By using the command prompt.
- By using [Terraform Cloud](https://www.terraform.io/cloud), which – among other things – provides a graphical user interface to Terraform.
- By using a CI/CD application (continuous integration/continuous delivery application such as [Jenkins](https://www.jenkins.io) or [Travis CI](https://travis-ci.org).

In this article, we'll focus on running Terraform from the command prompt. But before we do that, we need to talk about how Terraform executes files.

## How Terraform Executes Files

We've already noted that Terraform files are plain-text files that use a .tf file extension. One thing that we *didn't* note is that, other than the file extensions, Terraform file names are completely arbitrary. Do you want to name your configuration file **akamai.tf**? Fine: then name the file akamai.tf. Would you rather name your file **create_security_configuration.tf**, or maybe **bob.tf**? That's fine, too: they're your files, and you can give them any name you want.

> **Note**. OK, but wouldn't a file name like **this_is_a_terraform_configuration_file_that_you_can_run_anytime_you_need_to_modify_a_security_policy.tf** be a bit tedious to type at the command prompt? Yes, it would. But, as we're about to see, Terraform typically doesn't require you to type filenames at the command prompt.

In other words, you can name your files anything you want. Likewise, and especially if you've put the Terraform executable in your PATH, you can store those files anywhere you want to. It's up to you.

Good question: we *did* say that you typically don't type in a filename when you run a Terraform command, didn't we? So then how does Terraform know which file you want to use? Believe it or not, it doesn't: instead, Terraform automatically executes all the **.tf** files it finds in the working folder.

What does that mean, and does it really matter? Well, suppose we have 2 .tf files in our working folder:

- **siem.tf**, which returns information about our SIEM (Security Information and Event Management) settings.
- **slow_post.tf**, which returns information about our slow post protection settings.

If we run the `terraform plan` command, we'll get back both our SIEM setting information and our slow post protection information:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/siem-and-slowpost.png)

And we got back both the SIEM settings and the slow post settings because Terraform automatically executed both of our .tf files. What if we had 100 .tf files in our working folder? Then Terraform would have dutifully executed all 100 of those files.

That sounds like a good thing and, more often than not, it is. However, this approach can also lead to problems. For example, suppose you had one Terraform configuration file that created a new security policy, and a second Terraform configuration file that deleted all your security policies. When you ran Terraform, hoping to create a new security policy, Terraform would do just that: the first Terraform configuration would create the new policy. When that was finished, however, the second Terraform configuration would run and would delete all your security policies. That's probably not what you had in mind.

Admittedly, that scenario might be a little far-fetched. But this one isn't. Suppose our SIEM settings file looks like this:

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
  name = var.security_configuration
}

data "akamai_appsec_siem_settings" "siem_settings" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "siem_settings_output" {
  value = data.akamai_appsec_siem_settings.siem_settings.output_text
}
```

And suppose our slow post protections file looks like this:

```
terraform {
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
  name = var.security_configuration
}
data "akamai_appsec_slow_post" "slow_post" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}
output "slow_post_output_text" {
  value = data.akamai_appsec_slow_post.slow_post.output_text
}
```

Individually, these two files run flawlessly. But look what happens if the two files are in the same folder and we try to run `terraform plan`:

```
│ Error: Duplicate required providers configuration
│
│   on slow_post.tf line 2, in terraform:
│    2:   required_providers {
│
│ A module may have only one required providers configuration. The required providers were previously configured at siem.tf:2,3-21.


│ Error: Duplicate provider configuration
│
│   on slow_post.tf line 9:
│    9: provider "akamai" {
│
│ A default (non-aliased) provider configuration for "akamai" was already given at siem.tf:9,1-18. If multiple configurations are required, set the
│ "alias" argument for alternative configurations


│ Error: Duplicate data "akamai_appsec_configuration" configuration
│
│   on slow_post.tf line 13:
│   13: data "akamai_appsec_configuration" "configuration" {
│
│ A akamai_appsec_configuration data resource named "configuration" was already declared at siem.tf:13,1-51. Resource names must be unique per type
│ in each module.
```

As you can see, we get a ton of duplication errors: duplicate required providers, duplicate providers, duplicate data. What's going on here?

Well, what's going on is this: we haven't given you the full story yet. We noted that Terraform runs all the .tf files in a folder, which implies that Terraform runs file 1, then runs file 2, then runs file 3, and so on. But that's not really how it works. Instead, Terraform effectively grabs all the .tf files in a folder, combines them into a *single* file, and then executes that one combined file. That's why we get duplication errors. For example, both our sample files use this block to declare the Akamai Terraform provider:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}
```

There's nothing wrong with that: you *have* to declare which provider you're working with. Unfortunately, when you combine the two files you end up with two provider declaration blocks:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }

terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
```

And that's simply not allowed in Terraform.

So how do you work around this issue? Well, one way is to store the files in separate folders; that way, you don't have to worry about duplication since, by default, Terraform only works with the .tf files it finds in the working folder. Another option would be to remove all the duplicated information from one of the files; for example, we might edit slow_post.tf so that it looks like this:

```
data "akamai_appsec_configuration" "configuration_2" {
  name = var.security_configuration
}

data "akamai_appsec_slow_post" "slow_post" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "slow_post_output_text" {
  value = data.akamai_appsec_slow_post.slow_post.output_text
}
```

After the naming collisions are removed, the two files should play nicely together:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/siem-and-slowpost.png)

This may or may not be a problem: if you never have more than one .tf file in a folder it *won't* be a problem. But if you do have multiple Terraform files in a single folder it's something you need to be aware of.

## Running Terraform from the Command Prompt

For the moment, let's assume you have a completed Terraform configuration (see this article for a closer look at Terraform configuration files). Now what?

Well, now comes the time to run the `terraform plan` and `terraform apply` commands. (Yes, there are plenty of other Terraform commands you can run but, for this introductory article, we aren't going to worry about them. If you'd like more information, take a peek at the official [Terraform command line tutorials](https://learn.hashicorp.com/collections/terraform/cli).) The `terraform plan` command provides a sort of dry run in which it does its best to ensure that your Terraform configuration is syntactically correct and that your Terraform infrastructure and your real-life infrastructure are in sync.

> **Note**. Yes, that's a somewhat-cursory explanation of what `terraform plan` does. But it's good enough for now.

By default, you call `terraform plan` just the way you'd expect to call it, typing the following from the command prompt:

```
terraform plan
```

In turn, Terraform checks the syntax of your Terraform configuration, verifies the infrastructure, and then creates a plan that details what it will do if you decide to implement this plan. For example, if everything passes muster, you'll see a response similar to this:

```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:

  + create

Terraform will perform the following actions:

  # akamai_appsec_bypass_network_lists.bypass_network_lists will be created

  + resource "akamai_appsec_bypass_network_lists" "bypass_network_lists" {
    + bypass_network_list = [
      + "107828_GMSNETWORK",
        ]
    + config_id           = 76478
    + id                  = (known after apply)
      }

Plan: 1 to add, 0 to change, 0 to destroy.
```

More often than not, the fact that no error was generated means that you can (in this example anyway) create a new network bypass list . Keep in mind, however, that this isn't always the case: that's because the `terraform plan` command doesn't do a full check of *everything* when generating a plan. For example, although the command says it's going to create a new bypass list, here's what could actually happen when you run the `terraform apply` command, the command that tries to implement the plan:

```
akamai_appsec_bypass_network_lists.bypass_network_lists: Creating...
╷
│ Error: Title: Invalid Input Error; Type: https://problems.luna.akamaiapis.net/appsec/error-types/INVALID-INPUT-ERROR; Detail: Config(76478) is not of type WAP (Web Application Protector).
│
│   with akamai_appsec_bypass_network_lists.bypass_network_lists,
│   on akamai.tf line 13, in resource "akamai_appsec_bypass_network_lists" "bypass_network_lists":
│   13: resource "akamai_appsec_bypass_network_lists" "bypass_network_lists" {
│
```

In this example, `terraform plan` didn't verify whether the underlying security configuration was configured as a Web Application Protector configuration. The `terraform apply` command *did* do this verification, however, which is why that command failed.

Admittedly, errors like that are relatively rare. Instead, `terraform plan` is more likely to report errors such as this:

```
│ Error: Invalid resource type
│
│   on akamai.tf line 16, in resource "akamai_appsec_penalty_obx" "penalty_box":
│   16: resource "akamai_appsec_penalty_obx" "penalty_box" {
│
│ The provider akamai/akamai does not support resource type "akamai_appsec_penalty_obx". Did you mean "akamai_appsec_penalty_box"?
```

In this example, we misspelled the resource name **akamai_appsec_penalty_box**. When we ran `terraform plan`, Terraform not only caught the error, but even made a suggestion as to what we *should* have entered as the resource name.

It's important to note that, if your plan does fail, you won't be able to run your configuration until the problems have all been fixed. If we try to ignore the response and run `terraform apply` anyway, that command fails with the exact same error:

```
│ Error: Invalid resource type
│
│   on akamai.tf line 16, in resource "akamai_appsec_penalty_obx" "penalty_box":
│   16: resource "akamai_appsec_penalty_obx" "penalty_box" {
│
│ The provider akamai/akamai does not support resource type "akamai_appsec_penalty_obx". Did you mean "akamai_appsec_penalty_box"?
```

As we noted a minute ago, `terraform apply` is the command that actually does things: it creates a new security policy, it updates your slow post protection settings, it returns information about all your rate policies. If you want to put your plan into action, all you have to do is enter the following at the command prompt:

```
terraform apply
```

When you call terraform apply, Terraform will verify your plan and then pause:

```
Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value:
```

To run your command, type **yes** (all lowercase letters) and then press ENTER. If you enter anything else (even **Yes**), your apply command is cancelled:

```
Apply cancelled.
```

If you *do* type **yes**, Terraform  runs your configuration and keeps you informed of its progress:

```
akamai_appsec_penalty_box.penalty_box: Creating...
akamai_appsec_penalty_box.penalty_box: Creation complete after 5s [id=76967]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

If you see the message **Apply complete!** well, needless to say, that's the message that you *want* to see.
