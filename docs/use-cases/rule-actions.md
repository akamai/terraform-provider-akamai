---
layout: "akamai"
page_title: "Modifying a Kona Rule Set Rule Action"
description: |-
  Modifying a Kona Rule Set Rule Action
---


# Modifying a Kona Rule Set Rule Action

Among other tools, Kona Site Defender uses a vast collection of common vulnerability and exposure (CVE) rules to help protect your website from specific attack. Each of these rules (collectively referred to as the Kona Rule Set or KRS) is designed to look for a specific exploit, and to take action (issue an alert, deny the request, take a custom course of action, or do nothing at all) anytime the rule is triggered. These rule actions are predefined by Akamai, but you can use Terraform to change the action assigned to any of your KRS rules. Do you feel that issuing an alert is not sufficient for a given set of circumstances? Would you prefer that requests be denied any time the rule is triggered? Then use Terraform to change the rule action from alert to deny.

In this documentation, we'll show you how to do just that.

## Viewing the Action Currently Assigned to a KRS Rule

To view the action currently assigned to a rule, use the [akamai_appsec_rules](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_rules) data source, being sure to specify the ID of the security policy and the ID of the rule you're interested in. For example, this simple API call returns the action assigned to rule **970002** and security policy **gms1_134637**:

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

## Modifying a Rule Action

To change the action assigned to a KSD rule, you use the [akamai_appsec_rule resource](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rule) and set the rule action to one of the following values:

- `alert`. Writes an entry to the log file any time a request triggers the rule.
- `deny`. Blocks the request using a predefined response.
- `deny_custom_{custom_deny_id}`. Blocks the request using a custom deny response that you create. Custom deny actions are discussed later in this documentation.
- `none`. Takes no action.

For example, the following Terraform configuration sets the rule action for the rule 970002 to **alert**:

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

Here we use the **akamai_appsec_rule** resource to change the rule action for our KRS rule. Which KRS rule? That's easy; the rule that:

- Resides in the security configuration we connected to.
- Is associated with the security policy **gms1_134637**.
- Has the rule ID **970002**.

When we run our configuration, we should get back output similar to this:

```
akamai_appsec_rule.rule: Creating...
akamai_appsec_rule.rule: Creation complete after 4s [id=58843:gms1_134637:970002]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

And if we rerun our original API call, we should see that the rule action has been changed to alert:

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

rule_action = "none"
```

That's all it takes.

## Working with Custom Denies

Custom denies provide a way for you to create a custom page or custom API response for rejected requests. These custom pages/responses serve at least two purposes:

- They help you maintain a positive and branded experience in case of a false positive result (e.g., the suspected web attack wasn't actually a web attack).
- They can be used to misdirect actual attackers away from your website.

We won't explain how to create custom denies here; see the documentation for the `akamai_appsec_custom_deny` resource for that information. Here, we'll simply show you how to retrieve a collection of your available custom denies, then show you to use one of those denies as a rule action.

## Viewing the Custom Denies Available for Use

To determine which custom denies are available to use, we can use a Terraform configuration similar to this one:

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

data "akamai_appsec_custom_deny" "custom_deny_list" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "custom_deny_list_output" {
  value = data.akamai_appsec_custom_deny.custom_deny_list.output_text
}
```

Again, there's nothing very complicated here: we connect to the **Documentation** security configuration, use the **akamai_appsec_custom_deny** data source to retrieve a collection of custom denies, then echo back the contents of that collection. As you can see, the Documentation configuration has a pair of custom denies:

```
+-------------------------------------+
| customDenyDS                        |
+-------------------+-----------------+
| ID                | NAME            |
+-------------------+-----------------+
| deny_custom_64386 | Operation       |
| deny_custom_68193 | new custom deny |
+-------------------+-----------------+
```

If we'd like more-detailed information about any one of these custom denies,  we can rerun this same configuration; the only thing we do different is reference the ID of the custom deny of interest (in this example, **deny_custom_64386**):

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
          id          = "deny_custom_64386"
          name        = "Operation"
          parameters  = [
             {
                  name  = "prevent_browser_cache"
                  value = "true"
                },
              {
                  name  = "response_body_content"
                  value = <<-EOT
                        %(AK_REFERENCE_ID)
                        <h1>
                        This is my custom Error Message
                        </h1>
                   EOT
                },
              {
                  name  = "response_content_type"
                  value = "text/html"
                },
              {
                  name  = "response_status_code"
                  value = "403"
                },
            ]
        }
     ]
   }
 )
```

To set the action for a rule to this custom deny we can use the exact same Terraform configuration we used when setting the rule action to **alert**. The one difference? Now set the `rule_action` argument to the ID of our custom deny:

```
resource "akamai_appsec_rule" "rule" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  rule_id            = 970002
  rule_action        = "deny_custom_64386"
 }
```
