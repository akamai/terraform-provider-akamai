---
layout: "akamai"
page_title: "Akamai: Activations"
subcategory: "Application Security"
description: |-
  Activations
---

# akamai_appsec_activations


The `akamai_appsec_activations` resource allows you to activate or deactivate a given security configuration version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_activations" "activation" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  network = "STAGING"
  notes  = "TEST Notes"
  notification_emails = [ "user@example.com" ]
}

```

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `notification_emails` - (Required) A bracketed, comma-separated list of email addresses that will be notified when the operation is complete.

* `network` - The network in which the security configuration should be activated. If supplied, must be either STAGING or PRODUCTION. If not supplied, STAGING will be assumed.

* `notes` - An optional text note describing this operation.

* `activate` - A boolean indicating whether to activate the specified configuration version. If not supplied, True is assumed.

## Attribute Reference

In addition to the arguments above, the following attribute is exported:

* `status` - The status of the operation. The following values are may be returned:

  * ACTIVATED
  * DEACTIVATED
  * FAILED



