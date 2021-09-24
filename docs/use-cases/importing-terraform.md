---
layout: "akamai"
page_title: "Importing a Terraform Resource from One Security Configuration to Another"
description: |-
  Importing a Terraform Resource from One Security Configuration to Another
---


# Importing a Terraform Resource from One Security Configuration to Another

If you have multiple security configurations, it's possible that you have elements (security policies, configuration settings, rate policies, etc.) in configuration A that you'd also like to have in configuration B; for example, you might like to have the same set of custom rules in both configurations. So how can you replicate items found in configuration A in configuration B? Well, obviously you can simply recreate those items, from scratch, in configuration B: create new security policies with the same exact values, create new configurations settings with the same exact values, create new rate policies with the same exact values. That works, but that process can also be slow, tedious, and prone to errors.

If you're thinking, “There must be a better, easier way to do this,” then give yourself a pat on the back: there *is* a better, easier way to do. As it turns out, Terraform provides a way for you to export data from configuration A and then import that data into configuration B. The process isn't fully automated – you will have to do some manual editing here and there – but the approach is still faster and easier than manually creating new items in configuration B.

In this article, we'll show a simple example: we'll export the prefetch settings from configuration A (or, more correctly, the configuration with the ID **58843**) and then import those same settings into configuration B (i.e., the configuration with the ID **76478**). We'll focus on the prefetch settings simply, so we can keep our sample configurations short, sweet, and easy to follow. However, when it comes to exporting and importing items from a configuration you aren't limited just to prefetch settings; instead, you can export and import any (or all) of the following:

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

In this documentation, we'll use the [akamai_appsec_export_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_export_configuration) data source to grab the prefetch settings from configuration file 58843 and to save those values to a local file (**export.tf**). We'll then add some additional code to export.tf, and use Terraform's apply command to assign those settings to configuration 76478.

If you're familiar with the Terraform language, you might be aware of the **import** command, a command that allows you to import items from one configuration to another. So why don't we just use the import command to carry out this operation? Well, that's mainly because the import command doesn't actually import items from Configuration A into Configuration B. Instead, it imports items from Configuration A into your Terraform state file (e.g., **terraform.tfstate**). Because the imported data still needs some manual modification, you're left with two choices: directly manipulate the state file (not recommended), or copy the information from the Terraform state to a standalone configuration file. That's going to result in some additional work, and some additional chances to errors to occur.

Which is why we use the **akamai_appsec_export_configuration** data source instead.

## Before We Begin

We'll explain how the export/import process works momentarily. Before we do that, however, we should clarify exactly what that process is going to do. As noted, we have two security configurations: 58843 and 76478. Currently the prefetch settings for the two configurations looks like this

| Property             | 58843 | 76478 |
| -------------------- | ----- | ----- |
| enable_app_layer     | true  | false |
| all_extensions       | false | false |
| enable_rate_controls | true  | false |
| extensions           | mp4   |       |


What we want to do, in effect, is copy the setting values from 58843 and use those values to configure the settings for 76478. If we succeed, then the two configurations will have the exact same prefetch setting values. In other words:

| Property             | 58843 | 76478 |
| -------------------- | ----- | ----- |
| enable_app_layer     | true  | true  |
| all_extensions       | false | false |
| enable_rate_controls | true  | true  |
| extensions           | mp4   | mp4   |

And now we're ready to talk about how you do this.

## Exporting Data from a Security Configuration

To export data from a security configuration, use the **akamai_appsec_export_configuration** data source, taking care to specify exactly what it is you want to export. For example, a Terraform configuration that exports prefetch settings will look similar to this:

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

There's nothing particularly special about the configuration. It begins, as all our configurations do, by calling the Akamai Terraform provider and by providing our authentication credentials. The configuration then uses this block to connect to **Configuration A**, the security configuration that has the setting values we want to export:

```
data "akamai_appsec_configuration" "configuration" {
  name = "Configuration A"
}

The exporting itself takes place with this block:

data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version   = 9
  search    = ["AdvancedSettingsPrefetch"]
}
```

Here, we're simply calling the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_configuration) data source and connecting that data source to our security configuration (technically, to version **9** of our security configuration). We then use the `search` argument to specify the items we want to export; in this case that's only the prefetch setting values:

```
search = ["AdvancedSettingsPrefetch"]
```

What if we have additional items we want to export? That's fine: we can just add those items to the `search` list:

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

> **Note:** The syntax **$(path.module)/** simply indicates that we want to create the file in the same folder as the Terraform executable. But that's entirely up to you: if you'd rather save the file to a different folder, just replace **$(path.module)/** with the path to that folder.

After that, we tell Terraform to take all the exported settings and setting values and write them to this new file:

```
content = data.akamai_appsec_export_configuration.export.output_text
```

Believe it or not, that's all we have to do. When we run our configuration, we'll get a new file (**export.tf**) that looks something like this:

```
// terraform import akamai_appsec_advanced_settings_prefetch.akamai_appsec_advanced_settings_prefetch 58843
resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 58843
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```

This file consists of two parts. Part 1 is a commented-out Terraform import command (in a Terraform configuration file, **//** is used to indicate a comment):

```
// terraform import akamai_appsec_advanced_settings_prefetch.akamai_appsec_advanced_settings_prefetch 58843
```

Part 2 consists of the settings exported from configuration 58843:

```
resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 58843
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```

These are the values we want copied to the target configuration. Before we can do that, however, we need to manually make a change to export.tf.

## Modifying the Export File

When you export data from configuration 58843, the resulting data file is – not surprisingly – all about configuration 58843. That's why you'll occasionally see references to the configuration ID in the file:

```
config_id = 58843
```

Like we said, that's to be expected: after all, these are the settings and values for configuration 58843. However, these are also the settings and values we're going to import into a different security configuration. If we leave the file as-is, those settings and values will be imported into configuration 58843; in other words, we'd import them right back into the configuration we just exported them from. Needless to say, that's not what we want; instead, we want to import those settings into configuration 76478. That means we need to change all references to 58843 to 76478. For example:

```
config_id = 76478
```

In our sample file, there's only one place where we need to make a change. However, if we exported multiple items we'll likely have to make this change (or a similar change) in multiple places.

After updating the config_id value our export file looks like this:

```
// terraform import akamai_appsec_advanced_settings_prefetch.akamai_appsec_advanced_settings_prefetch 58843

resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 76478
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```

Because the export command doesn't create blocks for declaring the Akamai provider and for presenting your Akamai credentials, you'll also need to add this information to the beginning of the export.tf file:

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
```

Oh, and what about the commented-out import command. Well, like we said, it's commented out which means you can leave it in or take it out: it won't affect the Terraform configuration in any way. But if you prefer your Terraform configurations be as short and sweet as possible, then you can delete the comment, meaning that the final version of export.tf will look like this:

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

resource "akamai_appsec_advanced_settings_prefetch" "akamai_appsec_advanced_settings_prefetch" {
 config_id            = 76478
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = true
 extensions           = ["mp4"]
 }
```

## Importing the Exported Settings

We now have a Terraform configuration containing everything we need to apply the same prefetch settings found in configuration 58843 to configuration 76478. That means that we can use the `terraform plan` command to do a quick syntax check, and then use `terraform apply` to import settings. (Yes, just like we'd do if we were updating the settings from scratch, without copying the values found in configuration 58843.)

After running `terraform apply`, we can use the `akamai_appsec_advanced_settings_prefetch` data source to verify that our setting values have been updated:

```
+----------------------------------------------------------------------+
| advancedSettingsPrefetchDS                                           |
+------------------+---------------+----------------------+------------+
| ENABLE APP LAYER | ALL EXTENSION | ENABLE RATE CONTROLS | EXTENSIONS |
+------------------+---------------+----------------------+------------+
| true             | false         | true                 | mp4       |
+------------------+---------------+----------------------+------------+
```

## One Thing to Watch for When Working with Multiple .tf Files

When working in a folder that contains multiple .tf, (in our case, **akamai.tf** and **export.tf**) you might run into problems like this one:

```
│ Error: Duplicate required providers configuration
│
│   on export.tf line 2, in terraform:
│    2:   required_providers {
│
│ A module may have only one required providers configuration. The required providers were previously configured at akamai.tf:2,3-21.
```

In this case, the “duplicate providers” error occurs because we've used this block in two different files (akamai.tf and export.tf):

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}
```

Does it really matter that we reference the Akamai provider in two different .tf files? Well, if those two files reside in the same folder, yes, it does matter. That's because of the way Terraform processes .tf files. By default, when you run a Terraform command like `terraform plan` or `terraform apply`, Terraform runs all the .tf files it finds in the working folder. That's why you never have to tell Terraform which .tf file to run: it's going to run all of them.

However, Terraform isn't going to run these files one-by-one; that is, it isn't going to run akamai.tf and then run export.tf. Instead, it's effectively going to combine both of those files into a single file, and then run that combined file. And that's where the problem occurs: both akamai.tf and export.tf have blocks that call the Akamai provider, and you can't call the same provider twice in a single configuration. The net result? The “duplicate providers” errors.

So how do you get around this issue? Well, one obvious way is to save export.tf to a different folder; if you do that, then you won't have to worry about having multiple .tf files in the same folder. Another solution is to temporarily rename the file akamai.tf; for example, you might tack a new file extension on the end (e.g., .**temp**) meaning you'll now have these two files:

- akamai.tf.temp
- export.tf

With akamai.tf temporarily renamed, you now have just one .tf file. Problem solved.

Alternatively, you can comment out everything in akamai.tf; as you probably know, commented lines aren't executed, which means those lines won't conflict with the content in export. tf. You can quickly comment out an entire .tf file by making **/*** the first line in the file and ***/** the last line in the file; everything in between will be commented out. For example:

```
/*
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}
*/
```

Depending on the editor you use to create your Terraform configurations, you'll also have a visual cue telling you that the lines won't be executed.  

To restore functionality, just remove the two comment markers (**/*** and ***/**).
