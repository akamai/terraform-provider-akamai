---
layout: "akamai"
page_title: "Creating a Rate Policy"
description: |-
  Creating a Rate Policy
---


# Creating a Rate Policy

A “classic” way to take down a website is to simply overwhelm the site with requests, transmitting so many requests that the site exhausts itself trying to keep up. This might be done maliciously (with the intent of crashing your web servers) or it might done inadvertently: for example, if you announce some special offer to anyone who visits your site in the next hour you might get so many visitors that this legitimate traffic ends up bringing the site down.

In other words, and in some cases. it's possible to have too much of a good thing. Because of that, it's important that you monitor and moderate the number and rate of all the requests you receive. In the Akamai world, managing request rates is primarily done by employing a set of rate control policies. These policies revolve around two measures:

- **averageThreshold**. Measures the average number of requests recorded during a two-minute interval. The threshold value is the total number of requests divided by 2 minutes (120 seconds).

- **burstThreshold**. Measures the average number of requests recorded during a 5-second interval.


When configuring a rate policy, the average threshold should always be less than the burst threshold. Why? Well, the burst threshold often measures a brief flurry of activity that disappears as quickly as it appears. By contrast, the average threshold measures a much longer period of sustained activity. A sustained rate of activity obvious has the capability to create more problems that a brief and transient rate of activity.

Rate policies are also designed to trigger only when certain conditions are met. For example, a policy might be configured to fire only when a request results in a specified HTTP response code (e.g., a 404 or 500 error).

Terraform provides a way to quickly (and easily) create rate policies: this is done by specifying rate policy properties and property values in a JSON file, and then running a Terraform configuration that creates a new policy based on those values. After a policy has been created, you can then use an additional Terraform block to assign an action to the policy: issue an alert if a policy threshold has been breached; deny the request if a policy threshold has been breached; etc.

For example, the following Terraform configuration creates a new rate policy and assigns the policy to the **Documentation** security configuration. After the policy has been created, the configuration then assigns a pair of actions to the new policy:

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
  name = “Documentation”
}

resource "akamai_appsec_rate_policy" "rate_policy" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  rate_policy =  file("${path.module}/rate_policy.json")
}

output "rate_policy_id" {
  value = akamai_appsec_rate_policy.rate_policy.rate_policy_id
}

resource  "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  rate_policy_id     = akamai_appsec_rate_policy.rate_policy.rate_policy_id
  ipv4_action        = "alert"
  ipv6_action        = "alert"
}
```

As you can see, this Terraform configuration is similar to other Akamai Terraform configurations. After identifying the provider and providing our credentials, we connect to the security configuration:

```
data "akamai_appsec_configuration" "configuration" {
  name = “Documentation”
}
```

We then use this block and the [akamai_appsec_rate_policy](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_policy)  resource to create the policy:

```
resource "akamai_appsec_rate_policy" "rate_policy" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  rate_policy =  file("${path.module}/rate_policy.json")
}
```

If you're looking for all the configuration values for our new policy, you won't find them here; instead, those values are defined in a JSON file named **rate_policy.json**. (Incidentally, that file name is arbitrary: you can give your fine any name you want). This argument tells Terraform that the properties and property values for the new rate policy should be read from **rate_policy.json**:

```
rate_policy = file("${path.module}/rate_policy.json")=
```

> **Note:** We'll look at a sample JSON file in a few minutes.

The syntax **${path.module}** is simply a shorthand way to specify that the JSON file is stored in the same folder as the Terraform executable. You don't have to store your JSON files in the same folder as the Terraform executable: just remember that, if you use a different folder, you'll need to specify the full path to the JSON file. Otherwise, Terraform won't be able to find it.

All that's left now is to echo back the ID of the new policy:

```
output "rate_policy_id" {
  value = akamai_appsec_rate_policy.rate_policy.rate_policy_id
}
```

Well, that and configure policy actions for the new policy.

## Configuring Rate Policy Actions

After you've created your rate policy, you'll want to configure the actions to be taken any time the policy is triggered (e.g., any time the **burstThreshold** is exceeded). There are four options available to you, and these options must be set for both IPv4 I and IPv6 IP addresses:

- **alert**. An alert is issued if the policy is triggered.
- **deny**. The request is denied if the policy is triggered.
It's recommended that you don't start out by setting a rate policy action to **deny**.
Instead, start by setting all your actions to **alert** and then spend a few days monitoring and fine-tuning your policy threshold before you begin denying requests.
If you do not do this, you run the risk of denying more requests than you really need to.
- **deny_custom_{custom_deny_id}**. Takes the action specified by the custom dey.
- **none**. No action of any kind is taken if the policy is triggered.

> **Note:** As you'll see later in this documentation, rate policies have a property named `sameActionOnIpv` that indicates whether the same action (for example, **deny**) is used on both IPv4 and IPv6 addresses. When setting a rate policy action, however, you must specify both the IPv4 and IPv6 actions. For example, if you don't include the IPv4 action, then your configuration will fail because a required argument (in this example, `ipv4_action`) is missing.

In our sample configuration, we use the following Terraform block to set the rate policy actions:

```
resource  "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  rate_policy_id     = akamai_appsec_rate_policy.rate_policy.rate_policy_id
  ipv4_action        = "alert"
  ipv6_action        = "alert"
}
```

To set the actions, we use the [akamai_appsec_rate_policy_action](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_rate_policy_action) resource, and specify the appropriate security configuration (`config_id`), security policy (`security_policy_id`), and our newly-created rate policy (`rate_policy_id`). To indicate the rate policy, we reference the ID of that policy:

```
rate_policy_id = akamai_appsec_rate_policy.rate_policy.rate_policy_id
```

At that point all that's left is to configure the IPv4 and IPv6 actions:

```
ipv4_action = "alert"
ipv6_action = "alert"
```

Keep in mind that these two actions don't have to be the same. Although it's recommended that new rate policies have both these actions set to **alert**, you can specify that IPv4 addresses trigger a different action than IPv6 addresses:

```
ipv4_action = "deny"
ipv6_action = "alert"
```

## The Rate Policy JSON File

A JSON file used to define rate policy properties and property values will look similar to this:

```
{
    "additionalMatchOptions": [{
        "positiveMatch": true,
        "type": "ResponseStatusCondition",
        "values": ["400", "401", "402", "403", "404", "405", "406", "407", "408", "409", "410", "500", "501", "502", "503", "504"]
    }],
    "averageThreshold": 5,
    "burstThreshold": 8,
    "clientIdentifier": "ip",
    "description": "An excessive error rate from the origin could indicate malicious activity by a bot scanning the site or a publishing error. In both cases, this would increase the origin traffic and could potentially destabilize it.",
    "matchType": "path",
    "name": "HTTP Response Codes",
    "pathMatchType": "Custom",
    "pathUriPositiveMatch": true,
    "requestType": "ForwardResponse",
    "sameActionOnIpv6": true,
    "type": "WAF",
    "useXForwardForHeaders": false
}
```

Although tThe rate policies that you create will use a JSON file *similar* to the one shown above, there will be differences depending on such things as your `matchType`, your `additionalMatchOptions`, etc. Rate policy properties available to you are briefly discussed in the following sections of the documentation.

#### Required Properties

Any rate policy JSON file you create must include the properties shown in the following table:

| Property         | Datatype | Description                                                  |
| ---------------- | -------- | ------------------------------------------------------------ |
| averageThreshold | integer  | Maximum number of allowed hits per second during any two-minute interval. |
| burstThreshold   | integer  | Maximum number of allowed hits per second during any five-second interval. |
| clientIdentifier | string   | Identifier used to identify and track request senders; this value is required only when using Web Application Firewall. Allowed values are:<br /><br />* **api-key**. Supported only for API match criteria.<br />* **ip-useragent**. Typically preferred over **ip**  when identifying a client.<br /><br />* **Ip**. Identifies clients by IP address.<br />* **cookie:value.**  Helps track requests over an individual session, even if the IP address changes. |
| matchType        | string   | Indicates the type of path matched by the policy allowed values are:<br /><br />* **path**. Matches website paths.<br />* **api**. Matches API paths. |
| name             | string   | Unique name assigned to a rate policy.                       |
| pathMatchType    | string   | Type of path to match in incoming requests. Allowed values are:<br /><br />* **AllRequests**. Matches an empty path or any path that ends in a trailing slash (/).<br />* **TopLevel** . Matches top-level hostnames only.<br />* **Custom**. Matches a specific path or path component.<br /><br />This property is only required when the **matchType** is set to **path**. |
| requestType      | string   | Type of request to count towards the rate policy's thresholds. Allowed values are:<br /><br />* **ClientRequest**. Counts client requests to edge servers.<br />* **ClientResponse**. Counts edge responses to the client.<br />* **ForwardResponse**. Counts origin responses to the client.<br />* **ForwardRequest**. Counts edge requests to your origin. |
| sameActionOnIpv6 | boolean  | Indicates whether the same rate policy action applies to both IPv6 traffic and IPv4 traffic. |
| type             | string   | Rate policy type. Allowed values are:<br /><br />* **WAF**. Web Application Firewall.<br />* **BOTMAN**. Bot Manager. |

#### Optional Properties

Optional rate policy properties are described in the following table:

| Property              | Datatype | Description                                                  |
| --------------------- | -------- | ------------------------------------------------------------ |
| description           | string   | Descriptive text about the policy.                           |
| hostnames             | array    | Array of hostnames that trigger a policy match. If a hostname is not in the array then that request will be ignored by the policy. |
| pathUriPositiveMatch  | boolean  | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). |
| useXForwardForHeaders | boolean  | Indicates whether the policy checks the contents of the **X-Forwarded-Fo**r header in incoming requests. |

#### The additionalMatchOptions Object

Specifies additional matching conditions for the rate policy. For example:

```
"additionalMatchOptions": [
  {
    "positiveMatch": false,
    "values": [
      "121989_DOCUMENTATION001",
      "060389_DOCUMENTATION002"
      ],
    "type": "NetworkListCondition"
  }
]
```

Properties of the `addtionalmatchOptions` object are described in the following table:

| Property      | Datatype | Description                                                  |
| ------------- | -------- | ------------------------------------------------------------ |
| properties    | string   | Match condition type. Allowed values are:<br /><br />* **IpAddressCondition**<br />* **NetworkListCondition**<br />* **RequestHeaderCondition**<br />* **RequestMethodCondition**<br />* **ResponseHeaderCondition**<br />* **ResponseStatusCondition**<br />* **UserAgentCondition**<br />* **AsNumberCondition**<br /><br />This value is required when using `additionalMatchOptions`. |
| positiveMatch | boolean  | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `additionalMatchOptions`. |
| values        | string   | List of values to match on. This value is required when using `additionalMatchOptions`. |

#### The apiSelectors Object

Specifies the API endpoints to match on. Note that thus object can only be used if the `matchType` is set to **api**. For example:

```
"apiSelectors": [
  {
    "apiDefinitionId": 602,
    "resourceIds": [
      748
      ],
    "undefinedResources": false,
    "definedResources": false
  }
]
```

Properties of the apiSelectors object are described in the following table:

| Property           | Datatype        | Description                                                  |
| ------------------ | --------------- | ------------------------------------------------------------ |
| apiDefinitionId    | integer         | Unique identifier of the API endpoint. This value is required when using `apiSelectors`. |
| resourceIds        | array (integer) | Unique identifiers of one or more API endpoint resources.    |
| undefinedResources | boolean         | If **true**, matches any resource not explicitly added to your API definition without having to include the resource ID in the `resourceIds` property . If **false**, matches only those undefined resources that are listed in the `resourceIds` property. |
| definedResources   | boolean         | If **true**, matches any resource explicitly added to your API definition without having to include the resource ID in the `resourceIds` property . If **false**, matches only those defined resources that are listed in the `resourceIds` property. |

#### The bodyParameters Object

Specifies the request body parameters to match on. For example:

```
"bodyParameters": [
  {
    "name": "Country",
    "values": [
      "US",
      "MX",
      "CA"
      ],
    "positiveMatch": true,
    "valueInRange": false
    }
  ]
```

Properties for the `bodyParameters` object are described in the following table:

| Property      | Datatype | Description                                                  |
| ------------- | -------- | ------------------------------------------------------------ |
| name          | string   | Name of the body parameter to match on. This value is required when using `bodyParameters`. |
| positiveMatch | boolean  | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `bodyParameters`. |
| valueInRange  | boolean  | When **true**, matches values inside the `values` range. (Note that your values must be specified as a range to use this property.) For example, if your value range us 2:6, any value between 2 and 6 (inclusive) is a match; values such as 1, 7, 9, or 14 do not match.<br /><br />When **false**. matches values that fall outside the specified range. |
| values        | string   | Body parameter values to match on. This value is required when using `bodyParameters`. |

#### The fileExtensions Object

Specifies the file extensions to match on. For example:

```
"fileExtensions": {
  "positiveMatch": false,
  "values": [
    "avi",
    "bmp",
    "jpg"
    ]
  }
```

Properties of the `fileExtensions` object are described in the following table:

| Property      | Datatype | Description                                                  |
| ------------- | -------- | ------------------------------------------------------------ |
| positiveMatch | boolean  | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `fileExtensions`. |
| values        | string   | List of file extensions to match on. This value is required when using `fileExtensions`. |

#### The path Object

Specifies the paths to match on. For example:

```
"path": {
  "positiveMatch": true,
  "values": [
    "/login/",
    "/user/"
  ]
}
```

Properties of the path object are described in the following table:

| Property      | Datatype | Description                                                  |
| ------------- | -------- | ------------------------------------------------------------ |
| positiveMatch | boolean  | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `path`. |
| values        | array    | List of paths to match on. This value is required when using `path`. |

#### The queryParameters Object

Specifies the query parameters to match on. For example:

```
"queryParameters": [
  {
    "name": "productId",
    "values": [
    "DOC_12",
    "DOC_11"
    ],
  "positiveMatch": true,
  "valueInRange": false
  }
]
```

Properties of the `queryParameters` object are described in the following table:

| Property      | Datatype | Description                                                  |
| ------------- | -------- | ------------------------------------------------------------ |
| name          | string   | Name of the query parameter to match on. This value is required when using `queryParameters`. |
| positiveMatch | boolean  | Indicates whether the policy is triggered if a match is found (**true**) or if a match is not found (**false**). This value is required when using `queryParameters`. |
| valueInRange  | boolean  | When **true**, matches values inside the `values` range. (Note that your values must be specified as a range to use this property.) For example, if your value range us 2:6, any value between 2 and 6 (inclusive) is a match; values such as 1, 7, 9, or 14 do not match.<br /><br />When **false**. matches values that fall outside the specified range. |
| values        | string   | List of query parameter values to match on. This value is required when using `queryParameters`. |

