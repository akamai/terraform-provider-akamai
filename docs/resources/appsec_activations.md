---
layout: "akamai"
page_title: "Akamai: Activations"
subcategory: "Application Security"
description: |-
  Activations
---

# akamai_appsec_activations

**Scopes**: Security configuration

Activates or deactivates a security configuration. Security configurations activated on the staging network can be used for testing and fine-tuning; security configurations activated on the production network are used to protect your actual websites.

**Related API Endpoint**: [/appsec/v1/activations](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postactivations)

## Example Usage

Basic usage:

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

resource "akamai_appsec_activations" "activation" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  network             = "STAGING"
  notes               = "This configuration was activated for testing purposes only."
  notification_emails = ["user@example.com"]
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration being activated.

- `notification_emails` (Required). JSON array containing the email addresses of the people to be notified when activation is complete.

- `network` (Optional). Network on which activation will occur; allowed values are:

  * **PRODUCTION**
  * **STAGING**

  If not included, activation takes place on the staging network.

- `notes` (Required). Brief description of the activation/deactivation process. Note that, if no attributes have changed since the last time you called the akamai_appsec_activations resource, neither activation nor deactivation takes place: that's because *something* must be different in order to trigger the activation/deactivation process. With that in mind, it's recommended that you always update the `notes` argument. That ensures that the resource will be called and that activation or deactivation will occur.

- `activate` (Optional). Set to **true** to activate the specified security configuration; set to **false** to deactivate the configuration. If not included, the security configuration will be activated.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `status`. Status of the operation. Valid values are:

  *	**ACTIVATED**
  *	**DEACTIVATED**
  *	**FAILED**

