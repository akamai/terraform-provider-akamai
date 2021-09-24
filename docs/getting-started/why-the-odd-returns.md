---
layout: "akamai"
page_title: "Why is My Terraform Plan Returning Information I Didn't Ask For?"
description: |-
  Why is My Terraform Plan Returning Information I Didn't Ask For?
---


# Why is My Terraform Plan Returning Information I Didn't Ask For?

Sometimes you can get unexpected results when running the `terraform plan` command. For example, suppose you have a single Terraform configuration file (i.e., a **.tf** file) that returns a list of all the security policies in a configuration. When you run `terraform plan`, however, you get back results similar to this:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/tmi.png)

As you can see, you *did* get back the list of security policies; however, you also got back information about your SIEM settings and your slow post protection settings. Where in the world did those SIEM settings and slow post protection settings even come from?

As it turns out, those unexpected settings are likely coming from your Terraform state file, a file used, among other things, to keep your Terraform infrastructure and your real-world infrastructure in sync. And keeping the two in sync is important: after all, when you return a list of security policies by using Terraform that list should match the security policies you see in Akamai Control Center. For example, on the left we see that Terraform has returned 15 policies while, on the right, Control Center has returned the same 1 policies:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/terraform-and-control-center.png)

That's the way it should be.

By default, Terraform uses a local file (**terraform.tfstate**) to maintain state information. Terraform.tfstate is simply a JSON file and, while we don't recommend it, you *can* open the file and view the contents. In the following screenshot, you can see both the SIEM and the slow post protection settings; that's because we recently ran a .tf file that returned that information:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/terraform-state-file.png)

Of course, having done that, we hasten to add that a better, and safer, way to see what's in your state file is to use the command `terraform state list`. That returns a list of all the resources currently in the state file:

```
data.akamai_appsec_configuration.configuration
data.akamai_appsec_configuration.configuration2
data.akamai_appsec_siem_settings.siem_settings
data.akamai_appsec_slow_post.slow_post
```

As a general rule, there's nothing harmful about having “extra” information in your state file. For example, if we were to use the `terraform apply` command to run ou configuration file (the one that returns security policy IDs), all we'll get back are those security policies: Terraform won't return the SIEM settings or the slow post protection settings. But that's if we run `terraform apply`. If we run `terraform plan`, well, that command doesn't clean up the state file the same way that `terraform apply` does. Because of that, extraneous data is sometimes returned when you call `terraform plan`. In turn, that can make it difficult to know what your plan is actually going to do.

So is there a way to clean up your state file? As it turns out, there is. For starters, run `terraform state list` to ifnromation about the items in your state file:

```
data.akamai_appsec_configuration.configuration
data.akamai_appsec_configuration.configuration2
data.akamai_appsec_siem_settings.siem_settings
data.akamai_appsec_slow_post.slow_post
```

Next, use the `terraform state rm` command to remove those items, one-by-one. For example:

```
terraform state rm data.akamai_appsec_configuration.configuration
```

After you've removed all the resources, run terraform state list again. This time you shouldn't get anything back, which means that your state file is clear:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/terraform-state-list.png)

Equally important, if you now run `terraform plan` you should only get back the list of security policy IDs:

![Terraform](https://techdocs.akamai.com/terraform/img/appsec/getting-started/policies-only.png)
