---
layout: "akamai"
page_title: "Creating a Match Target"
description: |-
  Creating a Match Target
---


# Creating a Match Target

Your website might get millions of requests each day, and you might have any number of security policies that help you handle those requests and help protect your site from potentially-malicious requests. Do you apply every single security policy to every single request? No. Instead, you use match targets to define which security policy (if any) should apply to a specific API, hostname, or path. Should a request come in that triggers a match target (for example, you might have a match target that scans for a specific set of file extensions), the security policy associated with the target goes into action, using protections such as rate controls, slow POST protections, and reputation controls to determine whether the request should be honored.

You can create match targets in Terraform by using a configuration similar to this:

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

resource "akamai_appsec_match_target" "match_target" {
  config_id    = data.akamai_appsec_configuration.configuration.config_id
  match_target =  file("${path.module}/match_targets.json")
}
```

In this configuration, we begin by defining **akamai** as our Terraform provider and by providing our authentication credentials. We then use this block to connect to the **Documentation** configuration:

```
data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
```

After the connection is made, we use the [akamai_appsec_match_target](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_match_target) resource and the following block to create the match target:

```
resource "akamai_appsec_match_target" "match_target" {
  config_id    = data.akamai_appsec_configuration.configuration.config_id
  match_target =  file("${path.module}/match_targets.json")
}
```

Only two things happen in the preceding block: we specify the ID of the configuration we want the new match target associated with (`config_id`) and we specify the path to the JSON file containing the match target properties and property settings. That's what this argument is for:

```
match_target =  file("${path.module}/match_targets.json"
```

In this example, we're using a JSON file named **match_targets.json**. That file name is arbitrary: it doesn't have to be named **match_targets.json**. Similarly, we've placed our JSON file in the same folder as the Terraform executable; that's what the syntax **${path.module}/** is for. This file path is also arbitrary: you can store your JSON file in any folder you want. Just be sure that you replace **${path.module}/** with the path to that folder.

## The Match Target JSON File

When you create a match target, the properties and property values for that target are typically defined in a JSON file. When you run your Terraform configuration, information is extracted from that file and used to configure the new match target. A JSON file for creating a match target looks similar to this:

```
{
    "type": "website",
    "isNegativePathMatch": false,
    "isNegativeFileExtensionMatch": false,
     "hostnames": [
        "akamai.com",
        "techdocs.akamai.com",
        "training.akamai.com"
    ],
    "fileExtensions": ["sfx", "py", "js", "jar", "html", "exe", "dll", "bat"],
    "securityPolicy": {
        "policyId": "gms1_134637"
    }
}
```

Keep in mind that your match target JSON files won't necessarily look exactly like the preceding file; that's because different match targets will have different sets of properties and property values. The following sections of this document provide information on the properties available for use when creating a match target.

## Required Arguments

The following argument must be included in all your match target JSON files:

| Argument | Datatype | Description                                                  |
| -------- | -------- | ------------------------------------------------------------ |
| type     | string   | Match target type. Allowed values are:<br /><br />* website<br />* api |

#### The securityPolicy Object

Specifies the security policy to be associated with the match target; this object is required in any match target JSON file you create. For example:

```
"securityPolicy": {
  "policyId": "gms1_134637"
  }
```

Arguments related to the security policy object are described in the following table:

| Argument | Datatype | Description                               |
| -------- | -------- | ----------------------------------------- |
| policyId | string   | Unique identifier of the security policy. |

## Optional Arguments

The arguments described in the following table are optional: they might (or might not) be required depending on the other arguments you include in your match target. For example, if your match target includes the `filePaths` or `fileExtensions` object then your JSON file *can't* include the `defaultFile` argument.

| Argument                     | Datatype | Description                                                  |
| ---------------------------- | -------- | ------------------------------------------------------------ |
| configId                     | integer  | Unique identifier of the security configuration containing the match target. |
| configVersion                | integer  | Version number of the security configuration associated with the match target. |
| defaultFile                  | string   | Specifies how path matching takes place. Allowed values are:<br /><br />* **NO_MATCH**. Excludes the default file from path matching.<br />* **BASE_MATCH**. Matches only requests for top-level hostnames that end  in a trailing slash.<br />* **RECURSIVE_MATCH**.  Matches all requests for paths that end in a trailing slash. |
| fileExtensions               | array    | File extensions that the match target scans for.             |
| filePaths                    | array    | File paths that the match target scans for.                  |
| hostnames                    | array    | Hostnames that the match target scans for.                   |
| isNegativeFileExtensionMatch | boolean  | If **true**, the match target is triggered if a match *isn't* found in the list of file extensions. |
| isNegativePathMatch          | boolean  | If **true**, the match target is triggered if a match *isn't* found in the list of file paths. |
| sequence                     | integer  | Ordinal position of the match target in the sequence of match targets. Match targets are processed in the specified order: the match target with the sequence value **1** is processed first, the match target with the sequence value **2** is processed second, etc. |

#### The apis Object

Specifies the API endpoints to match on. Note that argument can only be used if the match target's `type` is set to **api**.

Arguments associated with the apis object are described in the following table:

| Argument | Datatype | Description                            |
| -------- | -------- | -------------------------------------- |
| id       | integer  | Unique identifier of the API endpoint. |
| name     | string   | Name of the API endpoint name.         |

#### The byPassNetworkLists Object

The bypass network list provides a way for you to exempt one or more network lists from the Web Application Firewall. For example:

```
"bypassNetworkLists": [
  {
    "id": "1410_DOCUMENTATIONNETWORK",
    "name": "Documentation Network"
    }
]
```

Arguments associated with the bypassNetworkLists object are described in the following table:

| Argument | Datatype | Description                            |
| -------- | -------- | -------------------------------------- |
| id       | string   | Unique identifier of the network list. |
| name     | string   | Name of the network list.              |

## Match Target Sequence

By default, match targets are applied in the order in which they are created. For example, suppose you have two different match targets:

- One checks to see if there are any file extensions included in a list of file extensions.
- One that checks to see if the hostname is included in the specified list of hostnames.

Assuming that your match targets were created in the order shown above, then for each request:

1.	The request is examined to see if it includes any of the specified file extensions.
2.	If any of these file extensions are found, the request is examined to see if the hostname in on the list of hostnames.

Although that approach works, it might not be the most efficient route you can take. For example, suppose you have 3 hostnames on your list: hostnames A, B, and C. Further suppose that you get 1 million requests each day. That means that, 1 million times a day, you're checking a request for any of your specified file extensions and, if found, then checking to see if the hostname is on the list of hostnames. That's fine if the majority of your requests come from hostnames A, B, and C. But what if only a handful of requests come from those hosts? That means you're doing a detailed search for file extensions on 1 million requests, even though only 1,000 of those requests are coming from a targeted host.

In a case like that, you might be better off swapping the match target sequence order: start by looking at the hostname on each request instead of starting with the file extensions. After all, if the hostname isn't A, B, or C you're done: there's no need to check the file extensions associated with the request. Instead of checking file extensions on 1 million files, you're checking file extensions only on the 1,000 requests coming from hosts A, B, or C.

> **Note:** A good rule of thumb is to start off by applying your most general match targets first, and then start working down to the more-specific match targets. Is the target shape blue? If not, then it doesn't matter. If so, then start to whittle down to questions like is it a blue circle; is it a blue circle with white polka dots; are the polka dots less than 1” in diameter; are those small polka dots oval-shaped rather than circular; etc.

As we learned a moment ago, by default match targets are applied in the order they were created: the first match target you create is applied first, the second match target you create is applied second, and so on. So what if you want to change the order in which your match targets are applied? Fortunately, you can do that by running a Terraform configuration similar to this:

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

resource "akamai_appsec_match_target_sequence" "match_targets" {
  config_id             = data.akamai_appsec_configuration.configuration.config_id
  match_target_sequence =  file("${path.module}/match_targets_sequence.json")
}
```

All we're doing here is connecting to the **Documentation** security configuration and then calling the [akamai_appsec_match_target_sequence](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_match_target_sequence) resource. More specifically, we're telling this resource to set the match target order using information found in the file **match_targets_sequence.json**. That's what this argument does:

```
match_target_sequence =  file("${path.module}/match_targets_sequence.json")
```

In that line, the value f**ile("${path.module}/match_targets_sequence.json")** represents the path to the JSON file that contains the target ordering information. The syntax **${path.module}/** indicates that the JSON file (**match_targets_sequence.json**) is in the same folder as our Terraform executable file.

> **Note:** What if you want to put the JSON file in a different folder? That's fine: just be sure you include the full path to that folder in your Terraform configuration.

As you might have guessed, the key to reordering your match targets is to specify the desired target sequence in this JSON file.

### Determining the current match target sequence

For better or worse, Terraform doesn't return information about match target sequencing; the [akamai_appsec_match_targets](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_match_targets) data source returns only minimal information about your match targets:

```
+---------------------------------+
| matchTargetDS                   |
+---------+-------------+---------+
| ID      | POLICYID    | TYPE    |
+---------+-------------+---------+
| 3723387 | gms1_134637 | Website |
| 3722423 | gms1_134637 | Website |
| 3722616 | gms1_134637 | Website |
| 3722692 | gms1_134637 | Website |
| 3723385 | gms1_134637 | Website |
| 3723386 | 4321_106673 | Website |
| 3722626 | gms1_134637 | Website |
| 3722379 | gms1_134637 | Website |
+---------+-------------+---------+
```

However, you can use the Application Security API to return information about your match targets that includes a target's sequence number:

```
"securityPolicy": {
    "policyId": "gms1_134637"
    },
"sequence": 2,
"targetId": 3722423
```

### The Match Target JSON File

The JSON file used to sequence your match targets will look similar to this:

```
{
  "type": "website",
  "targetSequence": [
    {
      "targetId": 3722423,
      "sequence": 1
    },
    {
      "targetId": 2660693,
      "sequence": 2
    },
    {
      "targetId": 2712938,
      "sequence": 3
    },
    {
      "targetId": 2809154,
      "sequence": 4
    },
    {
      "targetId": 3023865,
      "sequence": 5
    },
    {
      "targetId": 3505726,
      "sequence": 6
    },
    {
      "targetId": 3722379,
      "sequence": 7
    }
  ]
}
```

This JSON file has two required properties: `type` (which specifies whether the sequencing is for website matches or api matches), and `targetSequence`, an object containing the `targetId` and `sequence` value for each of your match targets. Do you want match target 3722423 to be the first match target applied? Then set its sequence value to **1**:

    {
      "targetId": 3722423,
      "sequence": 1
    },

Simply continue in this fashion until you've configured all your match targets in the desired order.
