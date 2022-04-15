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

Note that activation fails if the security configuration includes one or more invalid hostnames. You can find these names in the resulting activation error message. To activate the configuration, remove the invalid hosts and try again.

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
  version             = 1
  network             = "STAGING"
  note                = "This configuration was activated for testing purposes only."
  notification_emails = ["user@example.com"]
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration being activated.
- `version` (Required). Version of the security configuration to be activated.
- `notification_emails` (Required). JSON array containing the email addresses of the people to be notified when activation is complete.
- `network` (Optional). Network on which activation will occur; allowed values are:
  * **PRODUCTION**
  * **STAGING**
  If not included, activation takes place on the staging network.
- `note` (Optional). Brief description of the activation/deactivation process.
   If not supplied, a default note will be generated using the current date and time.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `status`. Status of the operation. Valid values are:
  *	**ACTIVATED**
  *	**DEACTIVATED**
  *	**FAILED**
