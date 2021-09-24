---
layout: "akamai"
page_title: "Creating an Automated Attack Groups (AAG) Security Configuration"
description: |-
  Creating an Automated Attack Groups (AAG) Security Configuration
---


# Creating an Automated Attack Groups (AAG) Security Configuration

When Akamai released automated attack groups (AAG) in October 2018, the technology represented a revolutionary development in protecting websites against Internet attacks. Prior to AAG, attack management primarily revolved around rule management: downloading rules to protect websites from common vulnerabilities and exposures (CVE), configuring rules, updating rules, deleting obsolete – well, you get the idea. This approach worked: websites were better protected. But management was occasionally difficult and time-consuming.

Automated attack groups offered a much-more scalable, cloud-based approach to protecting websites. Rules still play an important role in AAG, but administrators no longer have to manage those rules one-by-one. Instead, rules are divided into a set of attack groups (for example, one group of rules helps protect sites from SQL injection attacks), and administrators only have to decide which attack groups they want to deploy. Likewise, there's no need to upload new rules or delete obsolete rules: instead, the rules included in an attack group are managed (e.g., uploaded and deleted) by Akamai.

This article explains how you can use Terraform to create a security configuration that leverages automated attack groups. To do this requires a Terraform configuration that can carry out the following steps:

1. Create an IP Network List
2. Activate the Network List
3. Create a Security Configuration
4. Create a Security Policy
5. Assign an IP Network List
6. Enable Network Protections
7. Create and Configure Rate Policies
8. Enable Rate Control Protections
9. Configure Logging Settings
10. Configure Prefetch Settings
11. Enable Slow Post Protections
12. Configure Slow Post Settings
13. Enable Web Application Firewall Protections
14. Configure the Web Application Firewall Mode
15. Configure Attack Group Settings
16. Enable and Configure the Penalty Box
17. Activate the Security Configuration

In this article, we'll explain how to use Terraform to carry out each one of these steps. However, we'll start by showing you a Terraform configuration that can carry out all 18 steps.

## The AAG Terraform Configuration

The Terraform configuration that creates our AAG security configuration is shown below. Admittedly, the configuration might look a bit intimidating at first; that's mainly because this one configuration is carrying out 17 separate actions. (Or even a few more, depending on what you want to count as a single action.) But don't worry: the configuration's bark is far worse than its bite. And, to prove that, in a minute we'll walk you through each of the blocks included in this configuration.

But first, here's a Terraform configuration that creates and enables an AAG security configuration:

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
  name        = "Documentation AAG Test Configuration"
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
    rate_policy =  file("${path.module}/rate_policy_1.json")
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
    rate_policy =  file("${path.module}/rate_policy_2.json")
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
  enabled = true
}

// Step 12: Configure Slow Post Settings
resource "akamai_appsec_slow_post" "slow_post" {
  config_id                  = akamai_appsec_configuration.create_config.config_id
  security_policy_id         = akamai_appsec_security_policy.security_policy_create.security_policy_id
  slow_rate_threshold_rate   = 10
  slow_rate_threshold_period = 30
  duration_threshold_timeout = 20
  slow_rate_action = "alert"
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

Now let's see if we can explain exactly what this configuration does, and how it goes about doing it.

## Step 1: Create an IP Network List

```
resource "akamai_networklist_network_list" "network_list" {
  name        = "Documentation Test Network"
  type        = "IP"
  description = "Network used for the AAG documentation example."
  list        = ["192.168.1.1","192.168.1.2","192.168.1.3","192.168.1.4"]
  mode        = "REPLACE"
}
```

Among other things, network lists provide a way for you (and your firewall) to manage clients based either on their IP address or on their geographic location. For example, if you want to block all IP addresses from 192.168.1.0 to 192.168.1.255 from passing through your firewall, you could set the list property of your network list to the Classless Inter-Domain Routing (CIDR) address **192.168.1.0/24**. Alternatively, you could block (or allow) all clients from Norway by setting the list property to **NO**, the ISO 3166 country code for Norway. Or you can do both: you can create as many network lists as you need.

To create a network list, use the [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list) resource and the following arguments:

| Argument    | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| name        | Name of the new network list. Names don't have to be unique: you can have multiple network lists that share the same name. However, when the list is created it will be issued a unique ID, a value comprised of a numeric prefix and the list name. (Or a variation of that name. For example, a list named **Documentation Network** will be given an ID similar to **108970_DOCUMENTATIONNETWORK**, with the blank space in the name being removed.) |
| type        | Indicates the type of addresses used on the list. Allowed values are:<br /><br />* **IP**. For IP/CIDR addresses.<br />* **GEO**. For ISO 3166 geographic codes.<br /><br />Note that you cannot mix IP/CIDR addresses and geographic codes on the same list. |
| description | Brief description of the network list.                       |
| list        | Array containing either the IP/CIDR addresses or the geographic codes to be added to the new network list. For example:<br /><br />`list = ["US", "CA", "MX"]`<br /><br />Note that the list value must be formatted as an array even if you are only adding a single item to that list:<br /><br />`list = ["US"]`<br /><br />Note, too that `list` is the one optional argument available to you: you don't have to include this argument in your configuration. However, leaving out the `list` argument also means that you'll be creating a network list that has no IP/CIDR addresses or geographic codes. |
| mode        | Set to **REPLACE** when creating a new network list.         |

## Step 2: Activate the Network List

```
resource "akamai_networklist_activations" "activation" {
  network_list_id     = akamai_networklist_network_list.network_list.uniqueid
  network             = "STAGING"
  notes               = "Activation of the AAG test network."
  notification_emails = ["gstemp@akamai.com","karim.nafir@mail.com"]
}
```

After a network list has been created, your firewall won't actually block (or allow) clients on that list until that list has been activated. You can read more about network activation in the article Managing Network Lists, In the meantime, network lists can be activated by using the [akamai_networklist_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_activations) resource and the following arguments:

| Argument            | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| network_list_id     | Unique identifier of the network list being activated. In our AAG Terraform configuration we refer to the network list ID like this:<br /><br />akamai_networklist_network_list.network_list.uniqueid<br /><br />That, as you might recall, is the ID assigned to the network list we created in the previous step. It probably goes without saying that we can't hardcode a network list ID in our sample configuration: after all, the ID won't exist until after we've called `terraform apply` and the network list has been created. |
| network             | Specifies the network that the network list is being activated for. Allowed values are:<br /><br />* **STAGING**. “Sandbox” network used for testing and fine-tuning. The staging network includes a small subset of Akamai edge servers but is not used to protect your actual website.<br />* **PRODUCTION**. Network lists activated on the production network are used to help protect your actual website.<br /><br />If this argument is omitted, the network list is automatically activated on the staging network |
| notes               | Arbitrary information about the network list and its activation. |
| notification_emails | JSON array of email addresses for the people who should be notified when the activation process finishes. |

## Step 3: Create a Security Configuration

```
resource "akamai_appsec_configuration" "create_config" {
  name        = "Documentation AAG Test Configuration"
  description = "This security configuration is used by the documentation team for testing purposes."
  contract_id = "1-3UW382"
  group_id    = 13139
  host_names  = ["llin.gsshappylearning.com"]
}
```

Security configurations are containers that house all the elements – security policies, rate policies, match targets, slow POST protection settings – that make up a website protection strategy. Before you can create any of these elements, you'll need a place for those elements to reside. That place is a security configuration.

As shown at the beginning of this step, security configurations are created by using the [akamai_appsec_configuration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_configuration) resource and the following arguments:

| Argument    | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| name        | Unique name to be assigned to the new configuration.         |
| description | Brief description of the configuration and its intended purpose. |
| contract_id | Akamai contract ID associated with the new configuration. You can use the akamai_appsec_contracts_groups data source to return information about  the contracts and groups available to you. |
| group_id    | Akamai group ID associated with the new configuration.       |
| host_names  | Names of the selectable hosts to be protected by the configuration. Note that names must be passed as an array; that's what the square brackets surrounding **"documentation.akamai.com"** are for. If you want to add multiple hostnames to the configuration, just separate the individual names by using commas. For example:<br /><br />`host_names = ["documentation.akamai.com", "training.akamai.com", "events.akamai.com"]`<br /><br />All security configurations must include at least one protected host. |


For more information about creating a security configuration, see the article Creating Security Configurations.

## Step 4: Create a Security Policy

```
resource "akamai_appsec_security_policy" "security_policy_create" {
  config_id              = akamai_appsec_configuration.create_config.config_id
  default_settings       = false
  security_policy_name   = "Documentation Security Policy"
  security_policy_prefix = "doc0"
}
```

Security policies are probably the single most important item found in a security configuration. And that shouldn't come as much of a surprise: after all, many of the other items used in a security configuration (rate policies, attack group settings, firewall allow and block lists, slow post protection settings, etc., etc.) must be associated with a security policy. Although you can create a security configuration without creating a security policy, a policy-less security configuration is of very little use.

Because of that, as soon as we create our security configuration our next step is to add a security policy to the configuration. (A security configuration can contain multiple security policies, but we'll create just one for now.) To create the policy, we'll use the [akamai_appsec_security_policy](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_security_policy) resource and the following arguments:

| Argument               | Description                                                  |
| ---------------------- | ------------------------------------------------------------ |
| config_id              | Unique identifier of the security configuration to be associated with the new policy. Note that, in our AAG Terraform configuration example, we always refer to the security configuration ID like this:<br /><br />akamai_appsec_configuration.create_config.config_id<br /><br />That value represents the ID of the security configuration created in the previous step. Needless to say, we can't hardcode the security configuration ID because that ID won't exist until we run `terraform apply` and create the configuration. |
| security_policy_name   | Unique name to be assigned to the new policy.                |
| security_policy_prefix | **Four-character prefix used to construct the security policy ID. For example, a policy with the ID gms1_134637** is composed of three parts:<br /><br />* The security policy prefix (**gms1**).<br />* An underscore (_).<br />* A random value supplied when the policy is created (**134637**). |
| default_settings       | If **true**, the policy is created using the default settings for a new security policy. If **false**, a “blank” security policy is created. In our sample configuration we'll set this value to **false**. |


For more information about creating security policies, see the article Creating a Security Policy.

## Step 5: Assign an IP Network List

```
resource "akamai_appsec_ip_geo" "akamai_appsec_ip_geo" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  mode               = "block"
  ip_network_lists   = [akamai_networklist_network_list.network_list.uniqueid]
}
```

As noted previously, network lists enable you to configure your firewall to automatically block (or allow) a set of clients based on either IP address or geographic location. After your list (or lists) have been created, you use the [akamai_appsec_ip_geo](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_ip_geo) resource to specify what you want done with these lists. The akamai_appsec_ip_geo resource accepts the following arguments:

| Argument                   | Description                                                  |
| -------------------------- | ------------------------------------------------------------ |
| config_id                  | Unique identifier of the security configuration associated with the network list. This is a required property. |
| security_policy_id         | Unique identifier of the security policy associated with the network lists. This is a required property. |
| mode                       | Indicates whether the networks that appear on either the geographic networks list or on the IP network list should be allowed to pass through the firewall. Valid values are:<br /><br />* **allow**. Only networks on the geographic/IP network lists are allowed through the firewall. All other networks are blocked.<br />* **block**. All networks are allowed through the firewall except for networks on the geographic/IP network lists. Clients on those networks are blocked.<br /><br />This is a required property. |
| geo_network_lists          | Geographic networks on this list are either allowed or blocked based on the value of the mode argument. |
| ip_network_lists           | Geographic networks on this list are either allowed or blocked based on the value of the mode argument. |
| exception_ip_network_lists | Networks on this list are always allowed through the firewall, regardless of the networks that do (or don't) appear on either the geographic or IP networks list. |


Note that, in our sample Terraform block, we use these two lines to block all the network lists associated with the security policy:

```
mode = "block"

ip_network_lists = [akamai_networklist_network_list.network_list.uniqueid]
```

This means that only clients that appear on the IP or the geographic network list will be prevented from going through the firewall.

## Step 6: Enable Network Protections

```
resource "akamai_appsec_ip_geo_protection" "protection" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  enabled            = true
}
```

After you've configured your network lists you'll need to enable network protections; if you don't, then your security policy won't enforce any of the settings applied to those lists. (For example, any lists you've designated for blocking won't actually be blocked.) To enable network protections, use the [akamai_appsec_ip_geo_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_ip_geo_protection) resource and set the enabled property to **true**.

Keep in mind that network protections are enabled on a security policy-by-security policy basis. If you have multiple security policies you'll need to enable network protections on each one.

## Step 7: Create Rate Policies

```
resource "akamai_appsec_rate_policy" "rate_policy_1" {
    config_id   = akamai_appsec_configuration.create_config.config_id
    rate_policy =  file("${path.module}/rate_policy_1.json")
}

resource "akamai_appsec_rate_policy_action" "rate_policy_actions_1" {
    config_id          = akamai_appsec_configuration.create_config.config_id
    security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
    rate_policy_id     =  akamai_appsec_rate_policy.rate_policy_1.id
    ipv4_action        = "deny"
    ipv6_action        = "deny"
 }
```

Rate policies help you monitor and moderate the number and  rate of all the requests you receive; in turn, this helps you prevent your website from being overwhelmed by a sudden deluge of requests (which could be either an attack of some kind or just an unexpected surge in legitimate traffic). You create rate policies by using the [akamai_appsec_rate_policy](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_policy) resource and the following arguments:

| Argument    | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| config_id   | Unique identifier of the security configuration associated with the rate policy. |
| rate_policy | File path to the JSON file containing configuration information for the rate policy. In our sample configuration, **$(path.module)/** indicates that the JSON file (**rate_policy_1.json**) is stored in the same folder as the Terraform executable. This isn't required: you can store your JSON files anywhere you want. Just make sure that you specify the full path so that Terraform can actually find those files. |

> **Note**. We won't delve into the JSON files in this article. See Creating a Rate Policy for more information about what one of these JSON files actually looks like.


That's all you need to do to create a rate policy. However, when you create a rate policy the rate policy action is automatically set to a null value; that means that nothing happens any time the policy is triggered. Because we'd prefer to have requests be denied if they trigger one of our rate policies, we create the rate policy and then use a Terraform block similar to this to set the action for that policy:

```
resource "akamai_appsec_rate_policy_action" "rate_policy_actions_1" {
    config_id          = akamai_appsec_configuration.create_config.config_id
    security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
    rate_policy_id     =  akamai_appsec_rate_policy.rate_policy_1.id
    ipv4_action        = "deny"
    ipv6_action        = "deny"
 }
```

As you can see, for this policy we've set both the IPv4 and the IPv6 actions to deny.

Note that, in our sample Terraform configuration, we create 3 rate policies; as a result, we've repeated the Terraform blocks shown at the beginning of this step 3 times. It's possible to use a for_each loop to repeat the same action multiple times in a single block of code; we'll show you an example of that when we configure attack group settings. In this case, however, we took the easy way out and simply used the same Terraform block over and over (changing only the variable names).

So why did we do that? Well, creating a rate policy action requires you to know, in advance, the rate policy ID. When trying to work in a for_each loop that's tricky; it's easy to encounter an error like this:

```
The "for_each" value depends on resource attributes that cannot be determined until apply, so Terraform cannot predict how many instances will be created.
```

To help you avoid that error, we took the easy way out. We'll cover the ins and outs of the for_each loop in a separate article.

## Step 8: Enable Rate Control Protections

```
resource "akamai_appsec_rate_protection" "protection" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  enabled            = true
}
```

After you've created your rate policies and assigned your rate policy actions, you'll need to enable rate control protections on your security policy; if you don't do this, then your rate control policies won't actually get used. Because rate control policies are enabled on the security policy (as opposed to having to enabled each and every policy), all you need to do is:

1.	Call the [akamai_appsec_rate_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_protection) resource.
2.	Specify the ID of your security configuration and your security policy.
3.	Set the `enabled` property to **true**.

## Step 9: Configure Logging Settings

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

The [akamai_appsec_advanced_settings_logging](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_advanced_settings_logging) resource enables you to specify which HTTP headers you want to log, and which ones you don't. You configure this information by using a JSON file similar to the following:

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

| Argument        | Description                                                  |
| --------------- | ------------------------------------------------------------ |
| allowSampling   | Set to **true** to enable HTTP header logging. Set to **false** to disable header logging. |
| cookies         | Specifies how, or even if, cookie headers (i.e., HTTP headers that reference cookies set by the server) should be logged. Allowed values are:<br /><br />* **all**. All cookie headers are logged.<br />* **none**. No cookie headers are logged.<br />* **exclude**. All cookie headers except the ones specified by the **type** argument are logged.<br />* **only**. Only the cookie headers specified by the **type** argument are allowed.<br /><br />For example:<br /><br />"cookies": {<br/>  "type": "exclude",<br/>  "values": [<br/>    "documentation-cookie=test”<br/>  ]<br/>} |
| standardHeaders | Specifies how, or even if, standard HTTP headers such as **User-Agent**, **Forwarded**, and **Referer**, should be logged. Allowed values are:<br /><br />* **all**. All standard headers are logged.<br />* **none**. No standard headers are logged.<br />* **exclude**. All standard headers except the ones specified by the **type** argument are logged.<br />* **only**. Only the standard headers specified by the **type** argument are allowed.<br /><br />For example:<br /><br />"standardHeaders": {<br/>  "type": "only",<br/>  "values": [<br/>    "User Agent", "Referer”<br/>  ]<br/>} |
| customHeaders   | Specifies how, or even if, custom headers (i.e., non-standard HTTP headers) should be logged. Allowed values are:<br /><br />* **all**. All custom headers are logged.<br />* **none**. No custom headers are logged.<br />* **exclude**. All custom headers except the ones specified by the **type** argument are logged.<br />* **only**. Only the custom headers specified by the **type** argument are allowed<br /><br />For example:<br /><br />"customHeaders": {<br/>  "type": "all"<br/>} |
| override        | If **true**, header data won't be logged for any for security events triggered by settings in the security configuration. |


To log HTTP headers in a security configuration you need to:

1.	Create a JSON file containing the logging criteria.
2.	Call the **akamai_appsec_advanced_settings_logging** resource, specifying the ID of the security configuration and the path to the JSON file.

In our sample Terraform block, we use this line to indicate the path to the JSON file:

```
logging = file("${path.module}/logging.json")
```

As we've seen elsewhere, the syntax **${path.module}/** indicates that the JSON file (logging.json) can be found in the same folder as the Terraform executable. This isn't a requirement: you can store the JSON file anywhere you want. Just be sure that, in your Terraform configuration, you include the full path to the file.

In our sample Terraform configuration, logging settings are applied to all the security policies to the security configuration. However, by including the optional `security_policy_id` argument we can apply these values to an individual policy. In a case like that, the logging settings applied to the policy take precedence over the logging settings applied to the security configuration.

## Step 10: Configure Prefetch Settings

```
resource "akamai_appsec_advanced_settings_prefetch" "prefetch" {
 config_id            = akamai_appsec_configuration.create_config.config_id
 enable_app_layer     = true
 all_extensions       = false
 enable_rate_controls = false
 extensions           = ["cgi","jsp","aspx","EMPTY_STRING","php","py","asp"]
}
```

By default, your Web Application Firewall only inspects external requests: requests that originate outside of your origin servers and Akamai's edge servers. Internal requests – that is, requests between your origin servers and Akamai's edge servers – typically aren't inspected, and typically don't need to be inspected. (As a general rule, these “prefetch” requests are safe, and inspecting each one doesn't do much besides slowing down your website.)

However, there might be times (for example, if you're concerned about prefetch-driven amplification attacks) when enabling prefetch is useful. If so, you can enable and configure your prefetch settings by using the [akamai_appsec_advanced_settings_prefetch](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_advanced_settings_prefetch) resource and the following arguments:

| Argument             | Description                                                  |
| -------------------- | ------------------------------------------------------------ |
| config_id            | Unique identifier of the security configuration associated with the prefetch settings. |
| enable_app_layer     | Set to **true** to enable prefetch request inspection.       |
| all_extensions       | Set to **true** to enable prefetch requestion inspections on all file extensions included in a request. To limit the file extensions to a specified set, set this value to **false** and then specify the target file extensions by using the extensions argument. |
| enable_rate_controls | Set to **true** to enable rate policy checking on prefetch requests. |
| extensions           | Specifies the file extensions that, when included in a request, trigger a prefetch inspection. Note that this argument should only be included when **all_extensions** is set to false. |


Prefetch settings apply to the entire security configuration.

## Step 11: Enable Slow POST Protections

```
resource "akamai_appsec_slowpost_protection" "protection" {
  config_id          = akamai_appsec_configuration.create_config.config_id  
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  enabled            = true
}
```

Denial of service (DOS) attacks are attacks in which a website is inundated with tons of requests, all sent in rapid-fire succession. DOS attacks are bad but, unfortunately, they aren't the only way to bring down a website: another common attack vector is to  slowly (*very* slowly) send a series of requests to a site. Because the requests, and the responses, take so long, the website spends its time waiting for the client to respond instead of spending its time handling requests from new (and legitimate) clients. To help guard against these slow POST attacks, use the [akamai_appsec_slowpost_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_slowpost_protection) resource and set the **enabled** property to **true**.

## Step 12: Configure Slow POST Settings

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

After slow POST protections have been enabled, you might want to adjust the slow POST configuration settings as well. If so, that can be done by using the [akamai_appsec_slow_post](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_slow_post) resource and the following arguments:

| Argument                   | Description                                                  |
| -------------------------- | ------------------------------------------------------------ |
| config_id                  | Unique identifier of the security configuration associated with the slow POST settings. |
| security_policy_id         | Unique identifier of the security policy associated with the slow POST settings. |
| slow_rate_threshold_rate   | Specifies the minimum rate (in bytes per second) that a request must achieve to avoid triggering the slow POST policy. The threshold rate represents the average number of bytes received during the slow rate threshold period. |
| slow_rate_threshold_period | Time period (in seconds) used to calculate the slow rate threshold rate. |
| duration_threshold_timeout | Specifies the maximum length of time (in seconds) that the server will wait for the first 8KB of a POST request body to be received. If the duration threshold expires before the request has completed or before the first 8KB have been received then the slow POST policy is triggered.<br /><br />Note that the duration threshold always takes precedence over the slow rate threshold. |
| slow_rate_action           | Specifies the action taken if the policy is triggered. Allowed values are:<br /><br />* **alert**. An alert is issued.<br />* **abort**. The request is abandoned. |

## Step 13: Enable Web Application Firewall Protections

```
resource "akamai_appsec_waf_protection" "akamai_appsec_waf_protection" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  enabled = true
}
```

In order to use the Web Application Firewall (WAF), that firewall must be enabled. To enable firewall protection, use the [akamai_appsec_waf_protection](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_waf_protection) resource and:

1.	Connect to the appropriate security configuration
2.	Connect to the appropriate security policy (WAF is enabled/disabled on a security policy-by-security policy basis).
3.	Set the **enabled** property to **true**.

But don't go just yet: after enabling the firewall you'll also want to configure the firewall mode and configure your attack group settings.

## Step 14: Configure the Web Application Firewall Mode

```
resource "akamai_appsec_waf_mode" "waf_mode" {
  config_id          = akamai_appsec_configuration.create_config.config_id
  security_policy_id = akamai_appsec_security_policy.security_policy_create.security_policy_id
  mode = "AAG"
}
```

The Web Application Firewall mode determines the way in which the rules in your Kona Rule Set are updated. When using automated attack groups, this value should be set to **AAG**: that ensures that Akamai will take care of updating the rules as needed. Setting the firewall mode to **KRS** puts the onus on you: you'll need to periodically, and manually, update the rules by yourself.

To specify the firewall mode, use the [akamai_appsec_waf_mode](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_waf_mode) resource and, when using attack groups, set the mode to AAG.

## Step 15: Configure Attack Group Settings

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

| Rule ID | Description                                                  |
| ------- | ------------------------------------------------------------ |
| 3000000 | A SQL injection attack consists of insertion or "injection" of a SQL query via the input data from the client to the application. A successful SQL injection exploit can read sensitive data from the database, modify database data (Insert/Update/Delete), execute administration operations on the database (such as shutdown the DBMS), recover the content of a given file present on the DBMS file system and in some cases issue commands to the operating system.<br /><br />One of the common ways to probe applications for SQL Injection vulnerabilities is to use the 'GROUP BY' and 'ORDER BY' clause. Prior to using these clause in a SQL statement, the hacker first terminates the current query's context (assuming user input is used in the WHERE clause), which could be either numeric or a string literal, and after the clauses and comments out the rest of the query. <br /><br />This rule triggers on HTTP requests, which contain SQL Injection probes that use the 'GROUP BY' and 'ORDER BY' clause as mentioned above, when they are sent as user-input. |


These rules have further been classified into the following categories:

- **SQL** (SQL Injection). Attack type in which malicious SQL queries are inserted into a data entry field and then executed. Execution of these queries often results in the attacker gaining access to personally-identifiable information about a website's users.
- **XSS** (Cross-Site Scripting). Attack type in which client-side scripts are added to a web page and thus made available to users, users who proceed to unwittingly execute those scripts.
- **CMD** (Command Injection). Attack type which enables arbitrary (and typically malicious) commands to be executed on a host's operating system.
- **HTTP** (HTPP Injection). Attack type in which malicious commands are included within the parameters of an HTTP request.
- **RFI** (Remote File Inclusion). Attack type in which a malefactor attempts to dynamically insert malicious code into an application.
- **PHP**. PHP Injection. Attack type in which a malicious PHP script is uploaded to a website. This often takes place by using a poorly-constructed upload form.
- **TROJAN**. Attack type in which malicious code poses as a legitimate app, script, or link, and tricks users into downloading and executing the malware on their local device.
- **DDOS** (Direct Denial of Service). Attack type designed to bring down (or at least severely disrupt) a website. Typically, this is done by overwhelming the site with tens of thousands of spurious requests.
- **IN** (Inbound Anomaly). Specifies the anomaly score of an inbound request. In anomaly scoring, requests aren't judged by a single rule; instead, multiple rules – and the past historical accuracy of those rules – are used to determine whether or not a request is malicious.
- **OUT** (Outbound Anomaly). Specifies the anomaly score of an outbound request.

Does any of this really matter to you? Well, if you're using automated attack groups, then, yes, it really *does* matter. That's because, with automated attack groups, you don't manage individual rules; instead, you manage these categories (i.e., these attack groups). For example, to deny requests that violate a SQL Injection rule all you have to do is set the SQL attack group action to **deny**; in turn, any request that triggers the group is denied. You don't have to know which rules are in the group, you don't have to know if new rules have been added or obsolete rules have been deleted, you just have to know how to set the SQL attack group to **deny**. In the Terraform block shown above, we set the attack group actions for all our attack groups to **deny**, which means that we're going to deny any request that triggers the group.

Admittedly, the block for configuring the attack group settings is a bit more complicated than most of the other Terraform blocks in this article. Because of that, we'll take a few minutes to walk you through the code, step-by-step.

The block starts off by calling the [akamai_appsec_attack_group](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_attack_group) resource; that's pretty routine. After that first line, however, we encounter this:

```
for_each = toset(["SQL", "XSS", "CMD", "HTTP", "RFI", "PHP", "TROJAN", "DDOS", "IN", "OUT"])
```

What's going on here? Well, here we're up a **for_each** loop that enables us to configure each individual attack group, one at a time. The attack group IDs (SQL, XSS, CMD, etc.) are configured as an array, and the **toset** function is used to convert this array into a set of strings, the datatype required for use with a for_each loop.

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

Here we specify the ID of the attack group to be configured. However, we don't hardcode the ID; instead we use the syntax **each.value**. When we use a for_each loop, the first time through the loop the each.value property represents the first value in the loop; for us, that's **SQL**. That means that, the first time through the loop, we'll be configuring the SQL attack group. When that's done, we'll loop around, and each.value will now represents **XSS**, the second value in the for_each loop. That means we'll now configure the XSS attack group. This continues until we've looped through all the values included in the for_each loop.

And what exactly are we configuring? We're simply setting the `attack_group_action` for each attack group to **deny**. In other words, any request that triggers any attack group will be denied.

## Step 16: Enable and Configure the Penalty Box

```
resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id              = akamai_appsec_configuration.create_config.config_id
  security_policy_id     = akamai_appsec_security_policy.security_policy_create.security_policy_id
  penalty_box_protection = true
  penalty_box_action     = "deny"
}
```

In hockey (and in a few other sports), players who commit more-egregious fouls are removed from the game and sent to the “penalty box.” Players in the penalty box, at least for purposes of the game, don't exist: they can't participate until they've served their time; in addition, and because they can't participate, they're pretty much ignored. Eventually players in the penalty box are allowed back into the game, but if they commit another foul they're returned to the penalty box and the whole cycle starts over again.

The Akamai penalty box serves a similar function: if a request triggers an attack group, the offending client is sent to the penalty box for 10 minutes. That means that all requests from that client are ignored during that 10-minute period. When time is up, the client can resume making requests, but another violation will send the client back to the penalty box for 10 more minutes.

> **Note**. OK, yes: there's a bit more nuance to the penalty box than what we've described here, but this is explanation enough for our purposes.

To employ the penalty box (available only if you're using automated attack groups), use the [akamai_appsec_penalty_box](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_penalty_box) resource and the following two arguments:

- `penalty_box_protection`. Set to **true** to enable the penalty box, or set to **false** to disable the penalty box.

- `penalty_box_action`. Set to **deny** to deny all requests from clients in the penalty box, or set to **alert** to issue an alert any time a client in the penalty box makes a request.


Note that the 10-minute timeout period is not configurable.

## Step 17: Activate the Security Configuration

```
resource "akamai_appsec_activations" "new_activation" {
    config_id           = akamai_appsec_configuration.create_config.config_id
    network             = "STAGING"
    notes               = "Activates the Documentation AAG Test Configuration on the staging network."
    activate            = true
    notification_emails = ["gstemp@akamai.com","karim.nafir@mail.com"]
}
```

When you create a security configuration, that configuration is automatically set to inactive; that means that the security configuration that isn't actually analyzing and taking action on requests sent to your website. For that to happen, you need to employ the [akamai_appsec_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_activations) resource to activate the configuration. We won't go into the hows and whys of activating a security configuration in this documentation; for those details, see the Activating a Security Configuration article instead. Here we'll simply note that, when activating a configuration, you must specify the network  where the configuration will be active. The **akamai_appsec_activations** resource gives you two choices when picking a network:

- **STAGING**. Typically, you'll start by activating a configuration on the staging network (as in the sample block shown above). The staging network consists of a small number of Akamai edge servers, and provides a sandbox environment for testing and fine-tuning your configuration. Note that a configuration on the staging network is not working with your actual website and your actual website requests. Again, this network is for testing and fine-tuning, not for protecting your website.

- **PRODUCTION**. After you're satisfied with the performance of your security configuration, you can activate that configuration on the production network. Once activated there, the configuration will be working with your actual website and your actual website requests.


Without going into too much detail, we can activate our security configuration by using the akamai_appsec_activations resource and the following arguments:

| Argument            | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| config_id           | Unique identifier of the configuration being activated.      |
| network             | Specifies which network the security configuration is being activated on. Allowed values are:<br /><br />* **staging**<br />* **production** |
| notes               | Arbitrary notes regarding the network and its activation status. |
| activate            | If **true**, the specified network is activated; if **false**, the specified network is deactivated. Note that this property is optional: if omitted, `activate` is set to **true** and the network is activated. |
| notification_emails | JSON array of email addresses representing the people who'll receive a notification email when the activation process completes. |

