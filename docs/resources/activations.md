---
layout: "akamai"
page_title: "Akamai: Activations"
subcategory: "APPSEC"
description: |-
  Activations
---

# resource_akamai_appsec_activations


The `resource_akamai_appsec_activations` resource allows you to create or re-use Activationss.

If the Activations already exists it will be used instead of creating a new one.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}


resource "akamai_appsec_export_config" "appsecconfigedge" {
    name = "Akamai Tools"
}

resource "akamai_appsec_activations" "appsecactivations" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.version
    network = "STAGING"
    notes  = "TEST Notes"
    activate = true
    notification_emails = ["martin@akava.io"]
}

```

## Argument Reference

The following arguments are supported:
${SCHEMA_ARG_REF}

