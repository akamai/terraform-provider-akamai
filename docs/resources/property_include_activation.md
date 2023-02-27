---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_include_activation (Beta)

Use the `akamai_property_include_activation` resource to activate your include and make available to a property.
You can also modify the activation time out with the `AKAMAI_ACTIVATION_TIMEOUT` environment variable, providing time in minutes. The default time out is 30 minutes.

## Basic usage

```hcl
resource "akamai_property_include_activation" "my_example" {
  include_id    = "inc_X12345"
  contract_id   = "C-0N7RAC7"
  group_id      = "X112233"
  version       = 1
  network       = "STAGING"
  notify_emails = [
      "example@example.com",
      "example2@example.com"
  ]
}
```

The activation time out. Here, set at 120 minutes.

```shell
$ export AKAMAI_ACTIVATION_TIMEOUT=120
```

## Argument reference

This resource supports these arguments:

* `include_id` - (Required) An include's unique ID with the optional `inc_` prefix.
* `contract_id` - (Required) A contract's unique ID, including the optional `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the optional `grp_` prefix.
* `version` - (Required) The version of the include you want to activate.
* `network` - (Required) The network for which the activation will be performed.
* `notify_emails` - (Required) The list of email addresses to notify when the activation status changes.
* `note` - (Optional) A log message assigned to the activation request.
* `auto_acknowledge_rule_warnings` - (Optional) Automatically acknowledge all rule warnings for activation and continue.

## Attributes reference

This resource returns this attribute:

* `validations` - The validation information in JSON format.
