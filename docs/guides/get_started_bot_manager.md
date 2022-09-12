---
layout: akamai
subcategory: Bot Manager
---

# Bot Manager

Use our Bot Manager subprovider to identify, track, and respond to bot activity on your domain or in your app.

### Before you begin

* Understand the [basics of Terraform](https://learn.hashicorp.com/terraform?utm_source=terraform_io).
* Complete the steps in [Get started](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider) and [Set up your authentication](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/akamai_provider_auth). 
* Subscribe to Bot Manager.
* Use ≥v2.3.1 of our Terraform provider.
* Be familiar with [Bot Manager](https://techdocs.akamai.com/bot-manager/docs) terms and concepts.

## Using JSON to specify argument values

Unlike some Terraform subproviders, Bot Manager doesn’t use a lot of arguments. For example, suppose an item has three settings: `setting_1`, `setting_2`, and `setting_3`. In many subproviders, these settings, and their values, are expressed as separate arguments:

```
setting_1 = 100
setting_2 = 200
setting_3 = 300
```

With the Bot Manager, however, these settings are usually configured a JSON array, and that array is then used as the value for a single argument (e.g., `item_argument`). That’s why the examples used in this documentation often include code that looks like this:

```
  item_argument = <<-EOF
{
     "setting_1": 100,
     "setting_2": 200,
     "setting_3": 300
}
EOF
```

As you can see, this code snippet has a single argument: `item_argument`. In addition, the value of that argument is a JSON array containing our three settings and their values:

```
{
     "setting_1": 100,
     "setting_2": 200,
     "setting_3": 300
}
```

And what’s the purpose of `<<-EOF` and `EOF`? Well, `<<-EOF` simply tells Terraform that the subsequent lines of code contain the value to be assigned to `item_argument`. As a result, Terraform reads in each line, not stopping until it see the `EOF` delimiter, which tells Terraform that the input is complete. In our simple little example, that means that Terraform reads in the following lines of code and then assigns that input to `item_argument`:

```
{
     "setting_1": 100,
     "setting_2": 200,
     "setting_3": 300
}
```

Instead of hard-coding these values and settings in your Terraform configuration you can:

1.	Save the JSON array to a separate file. For example, you might save your custom category settings to a file named `custom-category.json`.

2.	Use syntax similar to the following to read the JSON file and assign the contents of that file to the `custom_bot_category` argument:
    `custom_bot_category = file("${path.module}/custom-category.json")``

In the preceding example, the syntax `${path.module}` is a shorthand way to indicate that the JSON file is stored in the same folder as the Terraform executable. Note that you don't _have_ to store your JSON files in the same folder as the Terraform executable: store the file wherever you want. Just remember that, if you use a different folder, you need to specify the full path to that file in your code. Otherwise, Terraform won't be able to find it.

## Workflow

Bot Manager has a number capabilities, and some of those capabilities depend on whether you’re running the standard version of Bot Manager or the Premier version of Bot Manager. (See the [Bot Manager documentation](https://techdocs.akamai.com/bot-manager/docs) for details.) Because Bot Manager can do so many things, there’s no way to walk you through an exhaustive set of Bot Manager workflows. However, we _can_ walk you through a typical workflow, one in which you:

1.	Create a custom bot category.
2.	Modify the action assigned to that custom category.
3.	Create a custom bot and assigning that bot to the new category.
4.	Move an Akamai-defined bot from its Akamai-defined category to the custom category.

### Create a custom category

Custom categories serve at least two purposes in Bot Manager. First,  custom-defined bots you create need to be placed in a custom category. You can’t add bots to an Akamai-defined category. Before you can create a custom-defined bot you need to create a custom category you can assign that bot to.

In addition to that, you can move an Akamai-defined bot out of its Akamai-defined category, but _only_ if you move that bot to a customer category. An Akamai-defined bot can’t be moved into a different Akamai-defined category. 

To create a custom category, use a Terraform configuration similar to this:

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

resource "akamai_botman_custom_bot_category" "custom_bot_category" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  custom_bot_category = file("${path.module}/custom_category.json")
}
```

Although that configuration might look a bit intimidating, its bark is actually far worse than its bite. For example, take the first block of code:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}
```

That’s boilerplate code for loading the Akamai Terraform provider. This code can be used exactly as-is in all your Terraform configurations.

That’s also true of the second block of code, which loads the `edgerc` file, a file containing your [authentication credentials](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/akamai_provider_auth):

```
provider "akamai" {
  edgerc = "~/.edgerc"
}
```

Finally, this block of code returns the ID of your security configuration:

```
data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
```

And what if you don’t _have_ a security configuration named `Documentation`? That’s fine: just replace `Documentation` with the name of that configuration. For example:

```
data "akamai_appsec_configuration" "configuration" {
  name = "NorthAmerica"
}
```

Incidentally, this block is optional: if you prefer, you can specify the ID of your security configuration (e.g., `76982`) when you create your custom category. The advantage of including the preceding block of code is that you only have to remember the security configuration's name, something that’s often easier to recall than the configuration ID.

That brings us to our final block of code, which uses the `akamai_botman_custom_bot_category` resource to create a custom bot category named `Vendor bots`:

```
resource "akamai_botman_custom_bot_category" "custom_bot_category" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  custom_bot_category = file("${path.module}/custom_category.json")
}
```

That code block requires the two arguments shown below:

| **Argument** | **Description** |
| --- | --- |
| `config_id`	| Unique identifier of the security configuration. In our sample code, we use the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_configuration) data source to return the ID for the `Documentation` configuration. If you already know the ID of the configuration you can skip that step and simply use that ID as the value of the `config_id` argument:<br><br>`config_id = 76982` |
| `custom_bot_category`	| JSON array containing settings and setting values for the new bot category. |

### Modify the action assigned to a custom category

When you create a custom bot category you don’t specify the action taken when the category is triggered. Instead, the category’s action is automatically set for you. If you want to assign a different action to the category use the `akamai_botman_custom_bot_category_action` resource and a Terraform configuration similar to this:

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

resource "akamai_botman_custom_bot_category_action" "custom_category_action" {
  config_id                  = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id         = "gms1_134637"
  category_id                = "2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8"
  custom_bot_category_action = file("${path.module}/action.json")
}
```

The first three blocks of code you should already be familiar with. The first block loads the Akamai Terraform provider and the second block points to the `edgerc` file, the file containing your authentication credentials. Meanwhile, the third block connects you to the `Documentation` security configuration and, not coincidentally, returns the ID of that configuration. Again, that step is entirely optional. If you prefer, leave it out and simply specify the ID of the security configuration in the code that modifies the category action.

The code that modifies the category action uses the following arguments:

| **Argument** | **Description** |
| ---      | ---          |
| `config_id`	| Unique identifier of the security configuration. In our sample code, we used the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_configuration) data source to return the ID for the Documentation configuration.  If you already know the ID of the configuration you can skip that step and simply use that ID as the value of the `config_id` argument:<br><br>`config_id = 76982`|
| `security_policy_id` | Unique identifier of the security policy associated with the custom bot category. Use the Application Security module’s [akamai_appsec_security_policy](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_security_policy) data source to return the IDs of your security policies. |
| `category_id`	| Unique identifier of the custom category whose action is being changed. Use the `akamai_botman_custom_bot_category` data source to return the IDs of your custom bot categories. |
| `custom_bot_category_action`	| JSON file containing the action to be talen when the bot category is triggered. For more information about these actions, see [Predefined actions for bot detections](https://techdocs.akamai.com/bot-manager/docs/predefined-actions-bot).|

### Create a custom bot

Although Bot Manager currently recognizes some 1,400 bots, that’s clearly not all the bots ever created. For example, you, one of your partners, or one of your vendors might have internal bots employed to help maintain your site. It’s unlikely that Akamai would know about, let alone have provided definitions for, bots like these. That’s one reason why you might need to create custom-defined bots.

To create a custom-defined bot, you need to specify the conditions that enable Bot Manager to identity the bot. For detailed information on how to configure these conditions, see [Update a custom-defined bot](https://techdocs.akamai.com/bot-manager/reference/put-custom-defined-bot).

Note that, when you create a custom bot, that bot must be assigned to a custom bot category. That simply means that, if you haven’t already done so, you need to create a custom category _before_ you create a custom bot.

The following Terraform configuration uses the `akamai_botman_custom_defined_bot` resource to create a custom bot named `vendor-bot`:

Here's the Terraform code:

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

resource "akamai_botman_custom_defined_bot" "custom_defined_bot" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  custom_defined_bot = file("${path.module}/custom_bot.json")
}

```
Like our previous examples, this configuration starts off by:

1.	Loading the Akamai Terraform provider.
2.	Retrieving authentication credentials from the `edgerc` file.
3.	Connecting to the `Documentation` security configuration.

Following those three steps (and those three blocks of code), the configuration calls the `akamai_botman_custom_defined_bot` resource and creates the new bot. The code for doing that includes these two arguments:

| **Argument**	| **Description** |
| --- | --- |
| `config_id`	| Unique identifier of the security configuration. In our sample code, we used the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_configuration) data source to return the ID for the `Documentation` configuration.  If you already know the ID of the configuration you can skip that step and simply use that ID as the value of the `config_id` argument:<br><br>`config_id = 76982` |
| `custom_defined_bot`	| JSON array containing the settings and setting values for the new bot.  |

### Recategorize an existing bot

When you subscribe to Bot Manager, you gain access to a large number of Akamai-defined bots, with each bot assigned to an Akamai-defined category. In most cases, those Akamai-defined categories will suit your needs just fine. However, it’s possible that you might have certain bots that you’d prefer to associate with a different category.

In a case like that, you need to “recategorize” the bot, which means that you need to move the bot to the custom category of your choice. (This must be a custom category: you can’t add bots to an Akamai-defined category.) To recategorize a bot, use the `akamai_botman_recategorized_akamai_defined_bot` resource and a Terraform configuration similar to the following:

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

resource "akamai_botman_recategorized_akamai_defined_bot" "recategorized_bot" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  bot_id      = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
  category_id = "2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8"
}
```

As usual, most of this configuration is boilerplate: the first two blocks load the Akamai provider and specify your authentication credentials, and the third block connects you to a security configuration named `Documentation`. After that, you use the following three arguments to specify the bot you want to move and where you want to move it to:

| **Argument**	| **Description** |
| --- | --- |
| `config_id`	| Unique identifier of the security configuration. In our sample code, we used the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_configuration) data source to return the ID for the `Documentation` configuration.  If you already know the ID of the configuration you can skip that step and simply use that ID as the value of the `config_id` argument:<br><br>`config_id = 76982` |
| `bot_id`	| Unique identifier of the Akamai-defined bot you want to move to a different category. Bot IDs can be returned by using the `akamai_botman_akamai_defined_bot` data source. |
| `category_id`	| Unique identifier of the custom category where the bot is being moved to. Use the `akamai_botman_custom_bot_category` data source to return the IDs of your custom bot categories. |
