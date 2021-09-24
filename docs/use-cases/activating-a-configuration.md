---
layout: "akamai"
page_title: "Activating a Security Configuration"
description: |-
  Activating a Security Configuration
---


# Activating a Security Configuration

By default, security configurations aren't activated when they're created; that simply means that those configurations aren't actually analyzing and responding to requests. Instead, and to actually make use of a security configuration, that configuration needs to be activated. Typically, that's a two-step process: first the configuration is activated on the staging network and then, after testing and fine-tuning, the configuration is activated on the production network. At that point, the configuration is fully deployed, and *is* analyzing and responding to requests.

If you're not familiar with the terms, the staging network is a relatively small set of Akamai Edge servers that simulates the conditions found on the production network (that is, conditions found in the real world). The staging network provides a sandbox environment where you can verify that your policies and configuration settings are working as expected, and can do so without putting your actual website at risk. For example, suppose you had, a rate policy that was effectively blocking all requests to your site. That policy could be discovered, and corrected, before it was put into production and before it started blocking all request to your actual site.

Just keep in mind that the one thing that the staging network doesn't do is performance testing. That's because, as noted, the staging network contains only a small number of Edge servers and a fraction of the server traffic found in the real world.

When the time comes, you can use Terraform to activate a configuration on either the staging network or the production network. A Terraform configuration for doing this will look similar to the following:

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
  activate            = true
  notes               = "This is a test configuration used by the documentation team."
  notification_emails = ["gstemp@akamai.com"]
}
```

For the most part, this is a typical Terraform configuration: we declare the Akamai provider, provide our authentication credentials, and connect to the **Documentation** configuration. We then use the [akamai_appsec_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_activations) resource and the following block to activate that configuration:

```
resource "akamai_appsec_activations" "activation" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  network             = "STAGING"
  activate            = true
  notes               = "This is a test configuration used by the documentation team."
  notification_emails = ["gstemp@akamai.com"]
}
```

Inside this block we need to include the following arguments and argument values:

| Argument            | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| config_id           | Unique identifier of the configuration being activated.      |
| network             | Name of the network the configuration is being activated for. Allowed values are:<br /><br />* STAGING<br />* PRODUCTION |
| activate            | If **true** (the default value), the security configuration will be activated; if **false**, the security configuration will be deactivated. Note that this argument is optional: if not included the security configuration will be activated. |
| notes               | Information about the configuration and its activation.      |
| notification_emails | JSON array of email addresses of the people who should be notified when activation is complete. To send notification emails to multiple people, separate the individual email addresses by using commas:<br /><br />notification_emails = ["gstemp@akamai.com", "karim.nafir@mail.com"] |


From here we can run `terraform plan` to verify our syntax, then run `terraform apply` to activate the security configuration. If everything goes as expected, you'll see output similar to the following:

```
akamai_appsec_activations.activation: Creating...
akamai_appsec_activations.activation: Creation complete after 2s [id=none]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## Reactivating a Security Configuration

Depending on the changes you make, a security configuration might need to be reactivated at some point. However, if you use the exact same Terraform block previously used in this documentation, reactivation won't take place. Instead, Terraform won't do anything at all:

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

As you can see, no resources were added, changed, or destroyed. In other words, nothing happened.

The problem here lies in the way that Terraform processes .tf files. When we originally activated the security configuration, we used the Terraform block shown above. When we tried to reactivate the configuration using that same block, Terraform was unable to see any changes: the `config_id` is the same, the `network` is the same, etc. Because nothing seems to have changed, Terraform did just that: nothing.

Perhaps the best way to work around this issue is to change the value assigned to the `notes` argument: this allows you to make a change of "some" kind without having to make a more drastic change (e.g., changing the network or the notification list). For example, when we originally activated the configuration, we used this line:

```
notes  = "This is a test configuration used by the documentation team."
```

To reactivate the security configuration, we simply need to change the value:

```
notes  = "This is a reactivated test configuration."
```

Now you should be able to reactivate the configuration. If you need to run activation again, just make a change to the `notes` argument.
