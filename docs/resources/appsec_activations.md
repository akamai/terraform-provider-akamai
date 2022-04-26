---
layout: "akamai"
page_title: "Akamai: Activations"
subcategory: "Application Security"
description: |-
  Activations
---

# akamai_appsec_activations

**Scopes**: Security configuration

Activates or deactivates a security configuration.
Security configurations activated on the staging network can be used for testing and fine-tuning; security configurations activated on the production network are used to protect your actual websites.
Note that activation fails if the security configuration includes one or more invalid hostnames. You can find these names in the resulting activation error message. To activate the configuration, remove the invalid hosts and try again.

### Important information if you upgrade to the 2.0.0 version of the Akamai Terraform provider

Prior to the release of version 2.0.0, the **akamai_appsec_activations** resource could either activate or deactivate a security configuration. If you are using version 2.0.0 or later the **akamai_appsec_activations** resource can no longer be used to deactivate configurations. In fact, the `activate` argument (which you could previously set to **true** or **false**) has been removed from the 2.0.0 version of the resource. Beginning with version 2.0.0, calling **akamai_appsec_activations** by using the `terraform apply` command automatically activates a security configuration.

To deactivate a security configuration create a Terraform configuration file that calls **akamai_appsec_activations** and references the configuration to be deactivated. However, instead of running `terraform apply`, use this command:

```
terraform destroy
```

When using the 2.0.0 provider you must also reference the version number of the security configuration being activated or, if you run `terraform destroy`, the version being deactivated. In previous versions of the Akamai provider, your Terraform configuration might include lines similar to these:

```
resource "akamai_appsec_activations" "activation" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  network             = "STAGING"
  notes               = "This configuration was activated for testing purposes only."
  notification_emails = ["user@example.com"]
}
```

Beginning with version 2.0.0, however, you must include the `version` argument and version number as well:

```
resource "akamai_appsec_activations" "activation" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  network             = "STAGING"
  note                = "This configuration was activated for testing purposes only."
  notification_emails = ["user@example.com"]
  version             = data.akamai_appsec_configuration.configuration.latest_version
}
```

In the preceding example, `version` is set to **latest_version**, a security configuration attribute that references the most recent version of the configuration. If you use **latest_version** (generally recommended) you’ll automatically activate the most recent configuration version. However, you can hard-code a specific version number if you prefer:

```
version = 5
```

Note that you do not have to upgrade to version 2.0.0. If you decide not to upgrade, the **akamai_appsec_activations** resource continues to function the way it has always functioned.


**Related API Endpoint**: [/appsec/v1/activations](https://techdocs.akamai.com/application-security/reference/post-activations)

## Example Usage

If you’re using a version of the Akamai Terraform provider released prior to version 2.0.0:

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
  activate            = true
}
```

If you’re using Akamai Terraform provider version 2.0.0 or later:

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
  note                = "This configuration was activated for testing purposes only."
  notification_emails = ["user@example.com"]
  version             = data.akamai_appsec_configuration.configuration.latest_version
}
```

## Argument Reference

The arguments available for the **akamai_appsec_activations** resource depend on which version of the Akamai Terraform provider you’re using. (You can determine your provider version by running the `terraform version` command.) If you are using a version of the Akamai provider released prior to version 2.0.0, these arguments are available to you:

- `config_id` (Required). Unique identifier of the security configuration being activated.

- `version` (Required). Version of the security configuration to be activated.

- `notification_emails` (Required). JSON array containing the email addresses of the people to be notified when activation is complete.

- `network` (Optional). Network on which activation will occur; if not included, activation takes place on the staging network. Allowed values are:
  * **PRODUCTION**
  * **STAGING**


- `notes` (required). Brief description of the activation or deactivation process. If no attributes have changed since the last time you called the **akamai_appsec_activations** resource, neither activation nor deactivation takes place. That's because something must be different in order to trigger one of these processes. Because of that, it's recommended that you always update the `notes` argument. Doing so ensures that the resource is called and activation or deactivation occurs. This argument applies only to versions prior to 2.0.0.

- `activate` (Optional). Set to **true** to activate the specified security configuration or set to **false** to deactivate the configuration. If not included, the security configuration is activated. This argument applies only to versions prior to 2.0.0.

If you’re running version 2.0.0 (or later), these arguments are available:

- `config_id` (Required). Unique identifier of the security configuration being activated. This is unchanged from previous versions.

- `notification_emails` (Required). JSON array containing the email addresses of the people to be notified when activation is complete. This is unchanged from previous versions.

- `network` (Optional). Network on which activation will occur; if not included, activation takes place on the staging network. Allowed values are:
    * **PRODUCTION**
    * **STAGING**


- `note` (Required). Brief description of the activation or deactivation process. If no attributes have changed since the last time you called the **akamai_appsec_activations** resource, neither activation nor deactivation takes place. That's because something must be different in order to trigger these processes. Because of that, it's recommended that you always update the **note** argument. That ensures that the resource is called and that activation or deactivation occurs.

    This argument (`note`, singular) is exactly the same as the `notes` argument (plural) used in previous versions. The only difference is that the name has changed.


  - `version` (Required). Version number of the security configuration being activated. This can be a hard-coded version number (for example, **5**), or you can use the security configuration’s **latest_version** attribute (data.akamai_appsec_configuration.configuration.latest_version). If you do the latter, you’ll always activate the most recent version of the configuration. This argument applies only to versions 2.0.0 and later.


## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `status`. Status of the operation. Valid values are:
  *	**ACTIVATED**
  *	**DEACTIVATED**
  *	**FAILED**
