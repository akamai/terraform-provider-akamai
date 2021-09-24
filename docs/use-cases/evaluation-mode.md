---
layout: "akamai"
page_title: "Working in Evaluation Mode"
description: |-
  Working in Evaluation Mode
---


# Working in Evaluation Mode

Evaluation mode provides a way for organizations to analyze the potential impact of changes to the Kona Rule Set (KRS) on their production network, but without having to immediately expose that network to the new ruleset. When evaluation mode is enabled, the latest version of the Kona Rule Set is downloaded and runs concurrently with the Kona Rule Set on the production network, but with one very important difference. The production network rules do what they've been configured to do: they issue an alert, they trigger a custom deny, they reject a request, etc. By comparison, the only thing that evaluation rules do is log the action that they *would* have taken had they been in force at the time, An alert rule won't actually issue an alert, and a deny rule won't actually deny a request: that remains the job of the KRS rules on the production network. But the evaluation rules *will* tell you, “Had I been in use, I would have denied this request.”

Admittedly, that might sound a little odd: what's the point of having a bunch of rules that don't do anything other than tell you what they *would* have done? Interestingly enough, however, that's the point: this allows you to analyze the effect of a new KRS ruleset on your actual network but without putting that network at risk. For example, suppose – based on your network and your configuration – a hypothetical Rule A would block every HTTP request sent to you website. (Obviously an unlikely scenario, but ….) Without evaluation mode, you might upgrade to the new ruleset only to discover that, suddenly, all requests to your site are being rejected. While you try to figure out what the problem is, and how you can fix it, users will be unable to connect to your site.

With evaluation mode, however, your site can continue to run as-is, fully protected by your current set of rules. In the meantime, you'll be able to monitor what the new ruleset would do if implemented. That means you'll be able to see the potential havoc that could be caused by Rule A before Rule A is actually implemented. That also means that you can make rule adjustments and fine-tune your site before upgrading to the latest KRS ruleset. That should help minimize any problems related to upgrading the ruleset, and help ensure that you can upgrade the rules (and take advantage of the new and improved protections) with little, if any, disruption to your production network.

Ideally, evaluation mode will help you answer key questions such as the following:

- Does the new ruleset recognize additional attack vectors, and do those attack vectors matter to me?
- Does the new ruleset cut down on the number of false positives generated in the site?
- Will we be able to re-enable rules that were disabled due to too many false positives?
- Are there rule exceptions that can be removed?

#### A Quick Note About How Evaluation Rules Works

Ignoring all the subtleties and complexities (at least for now), the basic idea of the Kona Rule Set is simple. An HTTP request – complete with its request line, its HTTP headers, and its message body –is received by a website and is analyzed by Kona Site Defender. For example, the request might include the **Connection** header and that header might include multiple instances of the **keep-alive** and the **close** options. In turn, Kona Site Defender checks its rule set to see if there are any rules involving the **Connection** header and the **keep-alive** and the **close** options. As it turns out, there's at least one such rule (rule **958295**):

> This rule inspects the 'Connection' header and looks for multiple values of the 'keep-alive' and 'close' options. Such behavior may be indicative of broken or malicious web clients, and may assist in detecting bot-generated traffic.
>
> This “triggers” rule 958295 and, in turn, your website might issue an alert or deny the request, depending on the action assigned to rule 958295.
>

If you're running in evaluation mode, each request received by a site is checked twice: you'll need to check the active KRS ruleset, and you'll also need to check the evaluation ruleset. (Remember, you're concurrently running two different rulesets.) On the surface, that means your website will need to do twice as much work, which will take twice as much time.

Fortunately, Akamai takes a number of steps to reduce this workload, primarily by ensuring that you don't necessarily have to query two different rulesets with each request. For example, suppose we have a production ruleset and an evaluation ruleset that look like this:

| Rule | Production Rule         | Evaluation Rule                             |
| ---- | ----------------------- | ------------------------------------------- |
| A    | Checks for condition A. | Checks for condition A.                     |
| B    | Checks for condition B. | Checks for condition B and for condition X. |
| C    | Checks for condition C. | Checks for condition C.                     |
| D    | Checks for condition D. | Checks for condition D.                     |
| E    | Checks for condition E. | Checks for condition E.                     |


As you can see, for 4 of our rules (rules **A**, **C**, **D**, and **E**) it doesn't matter whether you're using the production rule or the evaluation rule: the rules are identical. Consequently, there's no need to query both rulesets: if you query the production ruleset you'll get a result that matches what you'd get if you queried the evaluation ruleset. Because of that, one query will suffice: there's no need to query both rulesets.

With rule **B**, however, there are differences between the production version and the evaluation version. Because the rules differ, in this case you will need to query both rulesets.

The preceding scenario is highly-simplified, but the point remains the same: when running in evaluation mode, you only need to query the evaluation rules if those rules differ from the production rules. And because most of the rules in a ruleset remain unchanged when that set is upgraded, that greatly reduces the impact of running in evaluation mode.

## Enabling and Disabling Evaluation Mode

Evaluation mode is conducted on a single security configuration, and the process can be managed by using Terraform and the [akamai_appsec_eval](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_eval) resource. The operations (commands) available to you when calling this resource are summarized in the following table:

| Operation | Evaluation Value | Description                                                  |
| --------- | ---------------- | ------------------------------------------------------------ |
| START     | enabled          | Starts evaluation mode. By default, evaluation mode runs for four weeks, although you have the option of ending evaluation mode before the four-week period is up. |
| COMPLETE  | enabled          | Concludes the evaluation period (even if the four-week trial mode is not yet up) and automatically upgrades the Kona Rule Set on your production network to the same rule set you just finished evaluating. The **COMPLETE** command should be used when you have finished your evaluation and are sure that you want to upgrade your Kona Rule Set. If you want to temporarily pause evaluation (with the option of resuming that evaluation later), use the **STOP** command instead. |
| STOP      | disabled         | Pauses evaluation mode without upgrading the Kona Rule Set on your production network. The **STOP** command might be used if you want to pause evaluation mode, make some changes to your evaluation infrastructure, and then resume testing (to resume a stopped evaluation use the **RESTART** command). |
| RESTART   | enabled          | Resumes an evaluation trial that has been paused by using the **STOP** command. Note that, when you restart evaluation mode, the four-week time period begins where it left off; for example, if you stopped the trial after 3 weeks then you'll have one week of evaluation activities remaining after your issue the **RESTART** command. |
| UPDATE    | disabled         | Upgrades the Kona Rule Set rules in the evaluation ruleset to the latest version. Calling the **UPDATE** command does not update the same rules in your production network: it only upgrades the evaluation rules. |


A Terraform configuration for managing evaluation mode (in this example, for starting evaluation mode) will look similar to this:

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
resource "akamai_appsec_eval" "eval_operation" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = " "gms1_134637"
  eval_operation     = "START"
}

output "eval_mode_evaluating_ruleset" {
  value = akamai_appsec_eval.eval_operation.evaluating_ruleset
}

output "eval_mode_expiration_date" {
  value = akamai_appsec_eval.eval_operation.expiration_date
}

output "eval_mode_current_ruleset" {
  value = akamai_appsec_eval.eval_operation.current_ruleset
}

output "eval_mode_status" {
  value = akamai_appsec_eval.eval_operation.eval_status
}
```

## Managing Evaluation and Protected Hosts

To investigate the effects of a new ruleset on a specific host, you must first assign that host to the evaluation process: that can be either a host already protected by your security configuration, or a new host that hasn't been associated with a security configuration. One important thing to keep in mind is this: a host can be designated as an evaluation host, or it can be designated as a protected host, but it can't simultaneously be an evaluation host *and* a protected host. Does that matter? Well, it might. Suppose Host A is a protected host on your security configuration. If you assign Host A to the evaluation process then Host A *will no longer be a protected host*. That's not a good thing or a bad thing: it's just how the evaluation process works, and it's just something you need to keep in mind.

You can add one or more hostnames to an evaluation by using a Terraform configuration similar to this:

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

resource "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = ["documentation.akamai.com", "training.akamai.com", "videos.akamai.com"]
}
```

Most of this configuration file consists of boilerplate code you've already seen: we declare the Akamai Terraform provider, we provide the link to our authentication credentials (stored in the **.edgerc** file), and connect to the **Documentation** security configuration. From there we use the [akamai_appsec_eval_hostnames](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_eval_hostnames) resource to add hostnames to the evaluation infrastructure. That's something we do by specifying the configuration ID and by using this line:

```
hostnames = ["documentation.akamai.com", "training.akamai.com", "videos.akamai.com"]
```

All we're doing in the preceding line is assigning three hosts (**documentation.akamai.com**, **training.akamai.com**, and **videos.akamai.com**) to the `hostnames` property. If there's a “catch” here it's the fact that our set of hostnames must be specified as a JSON array; that's why the set of names is surrounded by a pair of square brackets (**[** ]). These square brackets are required even if we only specify a single host:

```
hostnames = ["documentation.akamai.com"]
```

And that's it: after we run `terraform apply` those three hosts will be used in the evaluation.

One thing we should point out here is that the hosts specified in the hostnames property *replace* any hosts currently configured for evaluation mode. For example, suppose we currently have a single host assigned to the hostnames property:

- test.akamai.com


After we run the configuration, the existing hostnames property values is erased, and is replaced by the set of host names specified in the configuration. In other words, now we have these three hosts assigned to the hostnames property:

- documentation.akamai.com
- training.akamai.com
- videos.akama.com

What about **test.akamai.com**? Well, that host was included in the configuration file; as a result, it's been deleted from the hostnames property. If we want to keep **test.akamai.com** and then add three additional hosts we have to specify all 4 hosts in the configuration file:

```
hostnames = ["documentation.akamai.com", "training.akamai.com", "videos.akamai.com", "test.akamai.com"]
```

If you want to remove all the hosts from evaluation mode simply set the value of the hostnames property to an empty array:

```
hostnames = []
```

> **Note**. Won't that cause problems if you remove all the hostnames from the evaluation list? No: you don't need to have any evaluation hosts if you don't want any. You might not get a lot of useful data without any evaluation hosts, but it won't cause any problems.

If you'd like to review the current list of evaluation hosts, just use the [akamai_appsec_eval_hostnames](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_failover_hostnames) data source and a Terraform configuration similar to the following:

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

data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "eval_hostnames" {
  value = data.akamai_appsec_eval_hostnames.eval_hostnames.hostnames
}
```

There are really only two things to point out here. First, we use this simple block to retrieve the collection of hostnames:

```
data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}
```

And then, when that's done, we use this block to display those hostnames:

```
output "eval_hostnames" {
  value = data.akamai_appsec_eval_hostnames.eval_hostnames.hostnames
}
```

## Protecting an Evaluated Host

As noted, evaluation hosts are not protected by the security configuration; that's the whole idea behind using evaluation hosts in the first place. (In fact, a protected host can't even be an evaluation host, at least not without removing the host from the security configuration's collection of protected hosts.) Not too surprisingly then, evaluation hosts are typically brand-new hosts that have never been under the protection of a production network.

As you evaluate your hosts, however, you might decide that it's worth moving them to the protected state (which also removes them from the evaluation infrastructure). If that's the case, you can use the [akamai_appsec_eval_protect_host](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_eval_protect_host) resource to easily move some (or even all) of your evaluation hosts to the protected hosts list. For example:

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

resource "akamai_appsec_eval_protect_host" "protect_host" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = ["training.akamai.com"]
}
```

In the preceding configuration, we connect to the **Documentation** security configuration, then use the **akamai_appsec_eval_protect_host** resource to protect the single host **training.akamai.com**, That's done by assigning the host to the `hostnames` property (as before, the hostnames value must be formatted as a JSON array):

```
hostnames = ["training.akamai.com"]
```

After this configuration runs, **training.akamai.com** will be a protected host, and we'll be left with just two evaluation hosts:

- documentation.akamai.com
- videos.aqkamai.com

Alternatively, you can use the following two Terraform blocks to protect *all* your evaluation hosts:

```
data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

resource "akamai_appsec_eval_protect_host" "protect_host" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = data.akamai_appsec_eval_hostnames.eval_hostnames.hostnames
}
```

In the first block, we use the **akamai_appsec_eval_hostnames** data source to return a collection of all our current evaluation hosts; that information gets stored in a variable named eval_hostnames. We then use that variable and the `hostnames` property when specifying the hosts to be moved to the protected state:

```
hostnames = data.akamai_appsec_eval_hostnames.eval_hostnames.hostnames
```

## Viewing Your Evaluation Rules

Evaluating a new Kona Rule Set doesn't mean that you enable evaluation mode, cross your fingers, and hope for the best. Instead, evaluation mode gives you an opportunity to modify and fine-tune that ruleset in order to come up with a configuration that maximizes protection while minimizing the problems such as false positives. While in evaluation mode, you can modify rule (and attack group) actions, conditions, and exceptions until you've created a configuration that best suits your unique infrastructure and security needs.

Of course, before your start modifying your evaluation rules you might want to determine which rules are included in your evaluation ruleset. (Evaluation rules are only the KRS rules that have been updated, which is typically a small subset of those rules.) To return information about your evaluation rules (which requires you to first start evaluation mode), use the [akamai_appsec_eval_rules](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_eval_rules) data source and a configuration file similar to this:

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

data "akamai_appsec_eval_rules" "eval_rule" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "output_text" {
  value = data.akamai_appsec_eval_rules.eval_rule.output_text
}
```

As you can see, there's not much to this configuration file: it kicks off by declaring the Akamai provider, indicates the path to our authentication credentials (the **.edgerc** file), and connects to the **Documentation** security configuration (the configuration being run in evaluation mode), From there we call the **akamai_appsec_eval_rules** data source specifying the configuration ID and the security policy ID. We then use this block of code to return the output from that query:

```
output "output_text" {
  value = data.akamai_appsec_eval_rules.eval_rule.output_text
}
```

That output will look similar to this:

```
+--------------------------------------------+
| RulesWithConditionExceptionDS              |
+---------+--------+------------+------------+
| ID      | ACTION | CONDITIONS | EXCEPTIONS |
+---------+--------+------------+------------+
| 699989  | alert  | False      | False      |
| 699990  | none   | False      | False      |
| 699991  | none   | False      | False      |
| 699994  | none   | False      | False      |
| 699995  | none   | False      | False      |
| 699996  | alert  | False      | False      |
| 950000  | alert  | False      | False      |
| 950001  | alert  | False      | False      |
```

In addition to that, you have other output options available to you:

| If you want to return …                            | … use a block similar to this                                |
| -------------------------------------------------- | ------------------------------------------------------------ |
| … just the rule actions                            | output "rule_action" {<br/>  value = data.akamai_appsec_rules.rule.rule_action<br/>} |
| … just the rule conditions and exceptions          | output "condition_exception" {<br/>  value = data.akamai_appsec_rules.rule.condition_exception<br/>} |
| … a JSON-formatted version of the rule information | output "json" {<br/>  value = data.akamai_appsec_rules.rule.json<br/>} |


If you only want information about a single rule, add the `rule_id` property and set the property value to the rule ID. For example, this block returns data only for rule **970002**:

```
data "akamai_appsec_eval_rules" "eval_rule" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  rule_id = 970002
}
```

## Modifying an Evaluation Rule

To modify a rule, use the [akamai_appsec_eval_rule](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_eval_rule) resource and a Terraform configuration like this:

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
resource "akamai_appsec_eval_rule" "eval_rule" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id  = "gms1_134637"
  rule_id             = 60029316
  rule_action         = "deny"
  condition_exception = file("${path.module}/condition_exception.json")
}
```

In this configuration, we start off the way most of our configurations start off: we declare, the Akamai provider, provide a pointer to our authentication credentials, and connect to the appropriate security configuration (in this case, the **Documentation** configuration). We then call the **akamai_appsec_eval_rule** resource, specifying the configuration, security policy, and rule IDs:

```
config_id          = data.akamai_appsec_configuration.configuration.config_id
security_policy_id = "gms1_134637"
rule_id            = 60029316
```

After that, we set these two properties of the rule:

- `rule_action`. Action that the rule takes when triggered. Allowed values are: **alert**; **deny**; **custom_deny**; and **none**. Keep in mind that, because these are evaluation rules, these are actually the actions that the rule *would* take when triggered. However, no action actually is taken.

  Note also that `rule_action` is a required property: it must be included in the configuration file even if you aren't changing the action (e.g., you're leaving the rule_action set to alert).

- `condition_exception`. Full path to a JSON file containing information about the rule's conditions and exceptions. Conditions are the criteria that specify when a rule is triggered; for example, a rule might fire if a requests uses a particular HTTP method or originates from a specified IP address or host domain, Exceptions are exactly what they sound like: exceptions that cause the rule *not* to be triggered. For example, you might say that the rule fires if the requests calls the POST method (a condition), but ignore that stipulation if the host happens to be on your internal network (an exception).


In our sample Terraform configuration, the property value is **file("${path.module}/condition_exception.json")**. In that value, **file** indicates that we are referencing a file path, and the syntax **$(path.module)/** indicates that the JSON file (named **condition_exception.json**) can be found in the same folder as the Terraform executable. Incidentally, this isn't required: the JSON file can be stored in any folder. Just be sure that you specify the full path to that folder and that file.

Note that `condition_exception` is optional: you can change the `rule_action` without having to change the conditions or exceptions.

And, remember, this only modifies the evaluation version of the rule. The production version of the rule, that rule actively protecting your production network, remains unchanged.

If you're running ASE, you also have the option of fine-tuning the individual attack groups used in evaluation mode. These attack groups include the following:

- **SQL** (SQL Injection). Attack type in which malicious SQL queries are inserted into a data entry field and then executed. Execution of these queries often result in the attacker gaining access to personally-identifiable information about a website's users.
- **XSS** (Cross-Site Scripting). Attack type in which client-side scripts are added to a web page and thus made available to users, users who proceed to unwittingly execute those scripts.
- **CMD** (Command Injection). Attack which enables arbitrary (and typically malicious) commands to be executed on a host's operating system.
- **HTTP** (HTTP Injection). Attack type in which malicious commands are included within the parameters of an HTTP request.
- **RFI** (Remote File Inclusion). Attack type in which a malefactor attempts to dynamically insert malicious code into an application.
- **PHP**. PHP Injection. Attack type in which a malicious PHP script is uploaded to a website. This often takes place by using a poorly-constructed upload form.
- **TROJAN**. Attack type in which malicious code poses as a legitimate app, script, or link, and tricks users into downloading and executing the malware on their local device.
- **DDOS** (Direct Denial of Service). Attack type designed to bring down (or at least severely disrupt) a website. Typically, this is accomplished by overwhelming the site with tens of thousands of spurious requests.
- **IN** (Inbound Anomaly). Specifies the anomaly score of an inbound request. In anomaly scoring, requests aren't judged by a single rule; instead, multiple rules – and the past historical accuracy of those rules – are used to determine whether or not a request is malicious.
- **OUT** (Outbound Anomaly). Specifies the anomaly score of an outbound request.

Similar to fine-tuning individual rules, you can modify the action and/or the conditions and exceptions used by an attack group. For example, the following Terraform configuration uses the [akamai_appsec_eval_group](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_eval_group) resource to modify the attack_group_action property for the SQL attack group. Note that the `condition_exception` property doesn't appear in this configuration. And that's fine: c`ondition_exception` is optional, which means that it *can* be omitted from a configuration.

Here's the sample file:

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
resource "akamai_appsec_eval_group" "eval_attack_group" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id  = "gms1_134637"
  attack_group        = "SQL"
  attack_group_action ="deny"
}
```

As you can see, this configuration file is very similar to the configuration file for updating an evaluation rule. That includes the following line, which sets the attack group action:

```
attack_group_action ="deny"
```

As with evaluation rules, the `action` property specifies what the rule would do if it was: 1) enabled on the production network; and, 2) triggered by a request. In this example, the `action` is set to **deny**, meaning that the request would be rejected.

And, like evaluation rules, this changes only affects the evaluation version of the attack group. The production network attack group is not modified in any way.
