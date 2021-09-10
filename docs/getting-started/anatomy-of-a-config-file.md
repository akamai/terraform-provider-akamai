---
layout: "akamai"
page_title: "The Anatomy of a Terraform Configuration File"
description: |-
  The Anatomy of a Terraform Configuration File
---


# The Anatomy of a Terraform Configuration File

As we've seen elsewhere, Terraform is based on the use of configuration files: plain-text files with .tf file extension. If you want to create a security policy or export data from a security configuration, you need to create – and then run – one of these configuration files. In principle, that's remarkably easy: create a file, run there file, start using your new security policy.

And if you're wondering, “OK, but what exactly does it mean to create a Terraform configuration file?” well, don't worry: in this documentation, we'll walk you through the creation of your first Terraform configuration file. That file will end up looking similar to this:

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

resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id     = var.security_policy_id
  penalty_box_protection = true
  penalty_box_action     = "deny"
}

output "list_contracts_and_groups" {
  value = data.akamai_appsec_contracts_groups.contracts_and_groups.output_text
}
```

Don't let this sample file scare you: it's nowhere near as complicated as it might look (at least in part because a lot of it is boilerplate text that can be copied from this sample and then pasted into *your* configuration file). And don't be dissuaded by the fact that it looks like computer programming: it's not. In fact, a Terraform configuration file is just that: it's a text file filled with configuration information. Yes, it might look a little like computer programming and, as you get more involved with Terraform you can incorporate some programming-like elements into your creations. For the most part, however, when you use Terraform all you're doing is saying, “Hey, I want to create this new thing, and I want the color to be red and the shape to be square.” Or, in our sample file, I want the **penalty_box_protection** property to be set to **true** and the **penalty_box_action** property to be set to **deny**:

```
penalty_box_protection = true
penalty_box_action = "deny"
```

In other words, even though there's an entire configuration language devoted to writing Terraform configuration files ([HashiCorp Language](https://www.terraform.io/docs/language/index.html), or HCL), those files don't need to be especially long or especially complicated. In fact, a fairly typical – but very useful – configuration only requires you to:

- Declare the Akamai Terraform provider.
- Provide your authentication credentials.
- Add the required data sources, resources, and outputs.

In this documentation, we'll explain exactly what we mean by all that.

## Declaring the Akamai Provider

When creating a .tf file, you start by telling Terraform which provider you want to work with. In our case, that's obviously going to be the Akamai provider. Consequently, all your Terraform configuration files should start off with the following block:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}
```

And, yes, there are other options that can be included in the **required_providers** block. However, because this is a getting started article, we're not going to worry about those options. You can read more about [provider options](https://www.terraform.io/docs/language/providers/index.html) on the Terraform documentation site.

One thing that we *will* mention, however, is the provider version. If you go to the home page for the Akamai provider, you'll see the latest version (**1.6.1** at the time this article was written) prominently displayed:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/akamai-provider.png)

That's important because Terraform doesn't automatically update your local setup any time a new provider version is released (in part because there could conceivably be times when you don't *want* to upgrade to the latest version). So how do you know which version of the Akamai provider you're actually using? One way to determine that is to run the following command from the command prompt:

```
terraform version
```

That command returns information about the version of Terraform you're running as well as the versions of any providers you've installed. As you can see, in this case we're a bit behind not only on our Terraform version but also on the Akamai provider version:

```
Terraform v1.0.2
on darwin_amd64

provider registry.terraform.io/akamai/akamai v1.6.0

provider registry.terraform.io/hashicorp/local v2.1.0

provider registry.terraform.io/hashicorp/time v0.7.2

Your version of Terraform is out of date! The latest version
is 1.0.4. You can update by downloading from https://www.terraform.io/downloads.html
```

As noted in the preceding response, you can upgrade Terraform itself by downloading the latest version from https://www.terraform.io/downloads.html and then replacing your existing (and outdated) Terraform executable with the new one. To upgrade the Akamai provider (and any other providers you might have installed) type the following command from the command prompt:

```
terraform init -upgrade
```

If everything goes well, you'll see output similar to the following:

```
Initializing the backend...

Initializing provider plugins...

Finding latest version of akamai/akamai...

Installing akamai/akamai v1.6.1...

Installed akamai/akamai v1.6.1 (signed by a HashiCorp partner, key ID A26ECDD8F0BCBA73)
```

And if we run terraform version a second time, we should see that our Akamai provider is now up-to-date (i.e., running version 1.6.1):

```
Terraform v1.0.2
on darwin_amd64

provider registry.terraform.io/akamai/akamai v1.6.1

provider registry.terraform.io/hashicorp/local v2.1.0

provider registry.terraform.io/hashicorp/time v0.7.2

Your version of Terraform is out of date! The latest version
is 1.0.4. You can update by downloading from https://www.terraform.io/downloads.html
```

## Providing Your Authentication Credentials

To manage your Akamai infrastructure by using Akamai Control Center, you need to have permission to carry out these management tasks: a random person off the street can't simply fire up Control Center and start messing around with your security configurations and rate policies. Not too surprisingly, these same permissions, and the same need to provide proof that you really *have* these permissions, is required in order to use the Terraform provider to manage your Akamai infrastructure.

That's why all the Terraform samples you'll see in our documentation includes this block:

```
provider "akamai" {
  edgerc = "~/.edgerc"
}
```

Before we go any further, two quick notes about this method of authentication (which involves storing your credentials in a file named **.edgerc**). 
First, this is not the only authentication method supported by Akamai; however, because it is the recommended method, it is the only one we are going to focus on.

Second, including this block in your .tf file is optional if – and this is an important if – you've created a file named **.edgerc** (note the dot at the beginning of the file name) and stored that file in your home folder. By default, the Akamai provider checks your home folder for a file named .edgerc; if the file exists, the provider then uses the credentials found in that file. Because this happens by default, you don't need to specify the credentials in your configuration file. We've included it simply because, in an introductory article, we don't want to take many shortcuts.

That leaves just one question: when we say that the .edgerc file should contain your “credentials,” what does that actually mean?

As it turns out, in order to use Akamai's Terraform provider you must first use Control Center and create an API client (a client, it goes without saying, that has the required permissions). We won't explain how to create an API client here; that's what this article's for. However, after you've created that client (and downloaded the credentials) you need to create an .edgerc file that looks similar to this one:

```
[default]
client_secret = /9kUbs/RFiAfidksRBk7DaXx5R9/cYTv/f132WDg5mQ=
host = akab-4rmgt2h36kujyc6d-jircmt2feb4lx5sk.luna.akamaiapis.net
access_token = akab-berscnnqusg3i2fq-4nru3b6yps4kmikr
client_token = akab-erqnjevsgkeir5ij-d4khuinb6oo4rejk
```

At that point, you're ready to start using Terraform.

## A Quick Note About Terraform Resources and Data Sources

Terraform configuration files consist of “blocks,” separate chunks of configuration code that carry out specific tasks. For example, the Terraform configuration shown at the beginning of this article includes the following two blocks:

```
data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  penalty_box_protection = true
  penalty_box_action = "deny"
}
```

As you probably noticed, the first block starts with the term **data** and the second block starts with the term **resource**. That's an important difference. In Terraform you'll primarily work with data sources and resources. For our immediate purposes, a data source returns data about some aspect of your Akamai infrastructure. For example, in the first block above we're using the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_configuration) data source to return information about one or more security configurations. The first “label” (**akamai_appsec_configuration**) referenced when calling the data source is the data source type, while the second label (**configuration**) is a variable (or, perhaps more accurately, a name) we can use when referring to the information returned by this data source, Note that the combination of type and name must be unique. Suppose our Terraform configuration includes these two blocks:

```
data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Training"
}
```

Because we've used the combination of **akamai_appsec_configuration** and **configuration** twice, we're going to get the following error:

```
│ Error: Duplicate data "akamai_appsec_configuration" configuration
│
│   on akamai.tf line 16:
│   16: data "akamai_appsec_configuration" "configuration" {
│
│ A akamai_appsec_configuration data resource named "configuration" was already declared at akamai.tf:13,1-51. Resource names must be unique per type
│ in each module.
```

If we need to call the **akamai_appsec_configuration** data source on multiple occasions, we'll have to change the name each time. For example, these two blocks of code will work, because the names are different:

```
data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_configuration" "configuration_2" {
  name = "Training"
}
```

After calling the data source we add a pair of curly braces ( **{ }** ) and define the instance of the data source (e.g., the name of the security configuration we want to return information for). “Defining an instance” means including arguments and argument values: in the expressions `name = "Documentation"`, **name** is an instance argument and **Documentation** is the argument value. Because name is a string, the argument value must be enclosed in double quotes. Those quotes can be left off if the argument is either a numeric value or a Boolean (true/false) value:

```
config_id = 58843
enabled = true
```

As you might expect, different data sources have different arguments. So then how are you supposed to know which arguments are available for which data source? That's easy; we just checked the [Akamai provider documentation](https://registry.terraform.io/providers/akamai/akamai/latest/docs) for the akamai_appsec_configuration data source:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/argument-reference.png)

> **Note**. And how did we know that akamai_appsec_configuration was the data source that returned configuration information in the first place? Again, the documentation is a good place to look for information like that.

Incidentally, you won't always need to use arguments, especially when dealing with data sources. For example, suppose you wanted to return information for all your configurations. In that case, you simply leave out the name argument:

```
data "akamai_appsec_configuration" "configuration" {
}
```

Resource blocks work the same wy as data source blocks work: the major difference is that you start the block off by calling, well, a resource:

```
resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  penalty_box_protection = true
  penalty_box_action = "deny"
}
```

Here, again you see the resource type (**akamai_appsec_penalty_box**) and the name (**penalty_box**), followed by a pair of curly braces and a set of arguments. (Although this isn't *always* the case, resources typically require more arguments than data sources.) In the preceding blocks, we have four arguments and argument values:

- `config_id = data.akamai_appsec_configuration.configuration.config_id`
- `security_policy_id = var.security_policy_id`
- `penalty_box_protection = true`
- `penalty_box_action = "deny`"

In the following table, we explain these arguments in a little more detail:

| Argument               | Value                                                      | Notes                                                        |
| ---------------------- | ---------------------------------------------------------- | ------------------------------------------------------------ |
| config_id              | data.akamai_appsec_ configuration.configuration. config_id | References the unique identifier of the Documentation security configuration. This is useful if you know the name of the configuration (**Documentation**) but aren't sure about the ID. If you *do* know the ID you can use syntax that references that ID instead:<br /><br />config_id = 58843<br /><br />In that case, you can also leave out the block that connects (by name) to the security configuration: seeing as how we know the ID, there's no need to connect by name. |
| security_policy_id     | var.security_policy_id                                     | The **var.** at the beginning of the argument value indicates that the security policy ID is being set to the value of a Terraform variable (in this case, a variable named **security_policy_id**). Variables provide a way for you to reference commonly-used values (such as security policy IDs or security configuration IDs) without having to remember the actual value assigned to those items: all you have to remember is the variable name.<br /><br />Variables also make it easy to switch from working with one set of items to another. Suppose you've been working exclusively with security configuration 58843; you'd now like to repeat those same tasks with security configuration 78643. One easy way to do that. Use a variable to represent the security configuration ID and then, as needed, simply change the value assigned to that variable.<br /><br />For more information on variables, including instructions on how to define those variables, see this article. |
| penalty_box_protection | true                                                       | The `penalty_box_protection` argument uses a Boolean datatype. Because of that, there's no need to put double quotes around the argument value. |
| penalty_box_action     | "deny"                                                     | The `penalty_box_action` argument uses the string datatype; hence the double quote marks surrounding the argument value. |

## Writing Your First Terraform Configuration File

> **Note**. If you have questions like “What should I name my Terraform configuration file?” or “Where should I store my Terraform configuration file?” see the Running Terraform Configuration Files. Here, we're going to assume you're already familiar with that information. We're also going to assume that you've created an API client and have added your credentials to an .edgerc file.

Let's now see if we can create a Terraform configuration file from scratch. To do that, open your favorite text editor and save a file with a **.tf** file extension (for example, **akamai.tf**) to the same folder where the Terraform executable is located. (As usual, you don't have to save your .tf files to that folder, but doing so makes the rest of this discussion less-complicated.) At the moment, your saved file should look something like this:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/blank-file.png)

In other words, it's an empty file.

If you remember back to the beginning of this article, we said that the first thing you needed to do when writing a Terraform configuration file was to declare the Akamai provider:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}
```

So let's start by adding that code to our new file:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/provider-only.png)

And even though it's optional, let's go ahead and add in our authentication information as well:

```
provider "akamai" {
  edgerc = "~/.edgerc"
}
```

That means our configuration file should now look like this:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/edgerc-included.png)

Let's now see which contracts and groups are associated with our account. (We chose this example because it will return data even if you haven't gotten around to creating your first security configuration.) By checking the provider documentation, we know that the [akamai_appsec_contracts_groups]() data source returns information about the Akamai contracts and groups associated with our account:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/contract-groups.png)

We also know that there are two arguments that can be used with this data source, both of which are optional:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/attributes-reference.png)

For our first try, we want to return information for all the contracts and all the groups. Because of that, we'll leave out the optional arguments:

```
data "akamai_appsec_contracts_groups" "contracts_and_groups" {
}
```

You'll notice that we:

- Started off with the word **data** to indicate that we're using a data source.
- Specified the data source **akamai_appsec_contracts_groups**.
- Included the data source name **contracts_and_groups**.
- Added an empty pair of curly braces (meaning that we aren't supplying any additional arguments).

That gives us a Terraform configuration that looks like this:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/data-scource-included.png)

We then add this block to indicate that we want to display the returned data in a table:

```
output "list_contracts_and_groups" {
  value = data.akamai_appsec_contracts_groups.contracts_and_groups.output_text
}
```

Before you ask, the value is composed of:

- The data source (**data.akamai_appsec_contracts_groups**).
- A period (**.**).
- The name we gave the data source (**contracts_and_groups**).
- Another period.
- One of the data source's output attributes.

And how did we know that `output_text` is a valid output attribute? Once again,all we had to do was take a peek at the documentation:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/output-reference.png)

As soon as we save our changes, we're ready to run the configuration. That step is covered in detail in the article Running Terraform (.tf) Files.
