---
layout: "akamai"
page_title: "Module: Network Lists"
description: |-
  Network Lists module for the Akamai Terraform Provider
---

# Network Lists Module Guide


## <a id="contents"></a>Table of contents

- [Retrieve network list information](#retrieve)
- [Create a network list](#create)
- [Activate a network list](#activate)
- [Modify an existing network list](#modify)
- [Deactivate a network list](#deactivate)
- [Import a network list](#import)


With the Akamai Network Lists provider for Terraform you can automate the creation, deployment, and management of network lists (shared sets of IP addresses, CIDR blocks, or broad geographic areas) used in various Akamai security products such as Kona Site Defender, Web App Protector, and Bot Manager. Along with managing your own lists, you can also access read-only lists that Akamai dynamically updates for you.

For more information about network lists, see the [API documentation](https://techdocs.akamai.com/network-lists/reference/api).


## Before you begin

This guide assumes that you have a basic understanding of Terraform and how it works (that is, you know how to install Terraform and the Akamai provider, how to configure your authentication credentials, how to create and use a .Terraform configuration file, etc.). If that’s not the case we strongly recommend you read the following two guides before going any further:

- [Akamai Provider: Get Started](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider)
- [Akamai Provider: Set Up Authentication](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/akamai_provider_auth)


## <a id="retrieve"></a>Retrieve network list information

You can obtain information about all the network lists available to you by using the [akamai_networklists_network_lists](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list) data source and its `output_text` attribute. To do that, use the following Terraform configuration:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_networklist_network_lists" "network_lists" {
}

output "network_lists_text" {
 value = data.akamai_networklist_network_lists.network_lists.output_text
}
```



## <a id="create"></a>Create a network list

As noted in the introduction, network lists are collections of device addresses that can be either:

- **Individual IP addresses or CIDR (Class Inter-Domain Routing) addresses**. CIDR provides a shortcut way to refer to a block of IP addresses. For example, the CIDR address 192.168.100.0/22 represent the 1,024 IPv4 addresses from 192.168.100.0 through 192.168.103.255.
- **Two-character ISO 3166 country codes** (**US** for the United States, **MX** for Mexico, etc.). Using country codes provides a way to manage network addresses across a broad geographical area.

These network lists can be used with Akamai products such as Kona Site Defender or Bot Manager, as well as by Terraform resources such as [akamai_appsec_ip_geo](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_ip_geo). That resource enables you to specify the network lists that are (or are not) allowed to pass through your IP/Geo firewall.

To create a new network list, use a Terraform configuration similar to the following:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

resource "akamai_networklist_network_list" "network_list" {
 name        = "Documentation Network"
 type        = "IP"
 description = "This is a test network used by the documentation team."
 list        = ["198.51.100.253","198.51.100.254","198.51.100.255"]
 mode        = "APPEND"
}

resource "akamai_networklist_activations" "activation" {
 network_list_id     = akamai_networklist_network_list.network_list.uniqueid
 network             = "STAGING"
 notes               = "This is the staged network for the documentation test network."
 notification_emails = ["gstemp@akamai.com"]
}
```

>  **Note**. The final block in this configuration activates the new network list. We'll ignore that block for now, then return to it in the **Activating a Network List** section of this documentation.

The Terraform configuration starts by initializing the Akamai provider and by providing the authentication credentials. After that, the configuration uses the [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list) resource and the following block to create the network list:

```
resource "akamai_networklist_network_list" "network_list" {
 name        = "Documentation Network"
 type        = "IP"
 description = "This is a test network used by the documentation team."
 list        = ["198.51.100.253","198.51.100.254","198.51.100.255"]
 mode        = "REPLACE"
}
```

The properties used in this block are described in the following table:

| **Argument**  | **Description**                                              |
| ------------- | ------------------------------------------------------------ |
| `name`        | Name of the new network list. Names don't have to be unique: you can have multiple network  lists that share the same name. When the list is created it’s issued a unique ID, a value comprised of a numeric prefix and the list name. In some cases, this might be a variation of that name. For example, a list named Documentation Network is given an ID similar to **108970_DOCUMENTATIONNETWORK**, with the blank space in the name being removed.) |
| `type`        | Indicates the type  of addresses used on the list. Allowed values are:     <br /><br />* **IP**. For IP/CIDR  addresses.  <br />* **GEO**. For ISO 3166  geographic codes.     <br /><br />Note that you can’t  mix IP/CIDR addresses and geographic codes on the same list. |
| `description` | Brief description of  the network list.                      |
| `list`        | Array containing  either the IP/CIDR addresses or the geographic codes to be added to the new  network list. For example:     <br /><br />`list =  ["US", "CA", "MX"]`     <br /><br />Note that the list  value must be formatted as an array even if you are only adding a single item  to that list:  <br /><br />`list =  ["US"]`     <br /><br />Note, too that `list` is the one optional argument available to you: you don't have to include this argument in your configuration. However, leaving out the `list` argument means that you'll end up creating a network list that has no IP/CIDR addresses or geographic codes (and, as a result, has little practical use). |
| `mode`        | Specifies whether  the addresses/geographic codes on the list should:     <br /><br />* Be added to the  existing set of addresses (**APPEND**).  <br />* Be removed from  the existing set of addresses (**REMOVE**).  <br />* Replace the  existing set of addresses (**REPLACE**).     <br /><br />Because we're creating a new network list, there won't *be* an existing set of  addresses. Consequently, you can set the mode either to **APPEND** or to **REPLACE**. |

After you've created your configuration, run the `terraform plan` command to validate and verify that configuration. Assuming your command succeeds you can then run `terraform apply` to create the network list. If the `apply` command is successful, you'll see output that includes the ID (**108972_DOCUMENTATIONNETWORK** ) of the network list:

```
akamai_networklist_network_list.network_list: Creation complete after 4s [id=108972_DOCUMENTATIONNETWORK]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```


## <a id="activate"></a>Activate a network list

After the network list is created, we use this Terraform block to activate the network:

```
resource "akamai_networklist_activations" "activation" {
 network_list_id     = akamai_networklist_network_list.network_list.uniqueid
 network             = "STAGING"
 notes               = "This is the staged network for the documentation test network."
 notification_emails = ["gstemp@akamai.com"]
}
```

This step is optional. You don't have to activate a network immediately after its creation. Instead, you can call the [akamai_networklist_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_activations) resource at any time to activate a network. Keep in mind, however, that a network can't be referenced by Akamai products or by other Terraform resources until it's been activated.

Our activation block uses the following arguments:

| **Argument**          | **Description**                                              |
| --------------------- | ------------------------------------------------------------ |
| `network_list_id`     | Unique identifier of  the network being activated. In this example,  **akamai_networklist_network_list.network_list.uniqueid** represents the ID of  the network that we just created. To activate a previously existing network, specify the ID of that network. |
| `network`             | Indicates which network the list will be activated on. Allowed values are:     <br /><br />* **STAGING**  <br />* **PRODUCTION**     <br /><br />Note that this  argument is optional. If not included in your configuration, the network is automatically activated on the staging network. |
| `notes`               | Optional description  of the network and information about the activation. |
| `notification_emails` | JSON array containing the email addresses of users to be notified when the network activation finishes. |

If you both create and activate a network list using a single Terraform configuration, you'll see output similar to the following:

```
akamai_networklist_network_list.network_list: Creating...
akamai_networklist_network_list.network_list: Creation complete after 4s [id=108970_DOCUMENTATIONNETWORK]
akamai_networklist_activations.activation: Creating...
akamai_networklist_activations.activation: Still creating... [10s elapsed]
akamai_networklist_activations.activation: Still creating... [20s elapsed]
akamai_networklist_activations.activation: Still creating... [30s elapsed]
akamai_networklist_activations.activation: Still creating... [40s elapsed]
akamai_networklist_activations.activation: Still creating... [50s elapsed]
akamai_networklist_activations.activation: Still creating... [1m0s elapsed]
akamai_networklist_activations.activation: Still creating... [1m10s elapsed]
akamai_networklist_activations.activation: Still creating... [1m20s elapsed]
akamai_networklist_activations.activation: Creation complete after 1m23s [id=6697234]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.
```

In addition, each user included in the `notification_emails` argument is emailed an activation notice.


## <a id="modify"></a>Modify an existing network list

Modifying an existing network list can be tricky, if only because you need to know which Terraform resource modifies which property values. Use the following table to help you determine which resource is the appropriate one to use:

| **Property**  | **Recommended resource for modifying the property**          |
| ------------- | ------------------------------------------------------------ |
| `name`        | [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list)     <br /><br />By  default, the network name is included in the network ID. For example, a  network named **USANetwork** will have an ID similar to **110987_USANETWORK**.  However, if you change the name of the network the network ID will not  change. You can change the name **USANetwork** to **HomeNetwork**, but  the ID for that network ID will still be **110987_USANETWORK**. |
| `type`        | [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list)     <br /><br />If  you change the network `type` (for example, if you switch from an IP list to a GEO list) you also need to change the list value as well.  |
| `description` | [akamai_networklist_description](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_description)    <br /><br /> You  can also use the [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list) resource to change the description. |
| `list`        | [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list)     <br /><br />Be sure to set the mode to the appropriate operation type (**APPEND**, **REPLACE**, **REMOVE**)  when changing the list. |
| `recipients`  | [akamai_networklist_subscription](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_subscription)     <br /><br />To  have email notifications sent to multiple addresses, separate those addresses by using commas:     <br /><br />["gstemp@akamai.com",  "karim.nafir@mail.com"] |

For example, this configuration changes the description assigned to the network list named **Documentation Network**:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_networklist_network_lists" "network_lists_filter" {
 name = "Documentation Network"
}

resource "akamai_networklist_description" "network_list_description" {
 network_list_id = data.akamai_networklist_network_lists.network_lists_filter.uniqueid
 name            = "Documentation Network"
 description     = "Updated description for the documentation network."
}
```

And this configuration changes the set of email addresses assigned to that same network list:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_networklist_network_lists" "network_lists_filter" {
 name = "Documentation Network"
}

resource "akamai_networklist_subscription" "subscribe" {
 network_list = data.akamai_networklist_network_lists.network_lists_filter.list
 recipients   = ["gstemp@akamai.com", "gmstemp@hotmail.com"]
}
```


##  <a id="deactivate"></a>Deactivate a network list

To deactivate a network list, remove all the addresses or geographic codes that have been assigned to that list. For example, this Terraform configuration deactivates the network list **Documentation Network**:

```
terraform {
 required_providers {
  akamai  = {
   source = "akamai/akamai"
  }
 }
}

provider "akamai" {
 edgerc = "~/.edgerc"
}

data "akamai_networklist_network_lists" "network_lists_filter" {
 name = "Documentation Network"
}

resource "akamai_networklist_network_list" "network_list" {
 name        = "Documentation Network"
 type        = "IP"
 description = "Test network list updated description."
 list        = []
 mode        = "REPLACE"
}
```

These two lines remove the addresses ore the geographic codes from the list:

```
list = []
mode = "REPLACE"
```

The first line sets the list (that is, the collection of device addresses or geographic codes) to an empty array. The second line then specifies that the new, empty list should replace the existing set of addresses/geographic codes. In effect, that removes all the IP addresses/geographic codes from the network list, which, in turn, deactivates that list.


## <a id="import"></a>Import a network list

Terraform allows you to add a resource to its state even if this resource was created outside of Terraform,
for example by using the Control Center application. This allows you to keep Terraform's state in sync with
the state of your actual infrastructure. To do this, use the `terraform import` command with a configuration
file that includes a description of the existing resource. The `import` command requires that you specify
both the `address` and `ID` of the resource. The `address` indicates the destination to which the resource
should be imported. This is formed by combining the resource type and local name of the resource as described
in the local configuration file, joining them with a period ("."). For example, suppose a network list has
been created outside of Terraform. You can use the information available in the Control Center to create a
matching description of this policy in your local configuration file. Here is an example, using sample values
for the resource's parameters:

```hcl
resource "akamai_networklist_network_list" "network_list" {
  name        = "Test-white-list"
  type        = "IP"
  description = "network list description"
  list        = ["198.51.100.0/24","203.0.113.255","233.252.0.0"]
  mode        = "APPEND"
}
```

The `address` of this resource is found by combining the resource type and its local name within the
configuration file: "akamai_networklist_network_list.network_list"

The `ID` indicates the unique identifier for this resource within Terraform's state. The unique identifier
for a network list can be found in the Control Center, and typically is of the form "12345_QWERTY". For this
example, suppose that the network list has been created and given the unique ID "80255_TESTWHITELIST". To
 import this resource into your local Terrform state, you would run this command:

```bash
$ terraform import akamai_networklist_network_list.network_list 80255_TESTWHITELISTNL
```


