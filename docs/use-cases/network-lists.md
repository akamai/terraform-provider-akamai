---
layout: "akamai"
page_title: "Managing Network Lists"
description: |-
  Managing Network Lists
---


# Managing Network Lists

Network lists are collections of device addresses; the addresses on a given list can be either:

- **Individual IP addresses or CIDR (Class Inter-Domain Routing) addresses**. CIDR provides a shortcut way to refer to a block of IP addresses; for example, the CIDR address 192.168.100.0/22 represents all the 1024 IPv4 addresses from 192.168.100.0 through 192.168.103.255.

- **Two-character ISO 3166 country codes** (for example, **US** for the United States, **MX** for Mexico, etc.). Using country codes provides a way for you to manage network addresses across a broad geographical area.


These network lists can then be used with Akamai products such as Kona Site Defender or Bot Manager, as well as by Terraform resources such as [akamai_appsec_ip_geo](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_ip_geo). That resource enables you to specify the network lists that are (or are not) allowed through your IP/Geo firewall.

This documentation covers such topics as:

- Creating a network list
- Modifying an existing network list
- Deactivating a network list

## Creating a Network List

To create a new network list, use a Terraform configuration similar to the following:

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

resource "akamai_networklist_network_list" "network_list" {
  name        = "Documentation Network"
  type        = "IP"
  description = "This is a test network used by the documentation team."
  list        = ["192.168.1.1","192.168.1.2","192.168.1.3"]
  mode        = "APPEND"
}

resource "akamai_networklist_activations" "activation" {
  network_list_id     = akamai_networklist_network_list.network_list.uniqueid
  network             = "STAGING"
  notes               = "This is the staged network for the documentation test network."
  notification_emails = ["gstemp@akamai.com"]
}
```

> **Note:** The final block in this configuration activates the new network list. We'll ignore that block for now, then return to it in the Activating a Network List section of this documentation.

The Terraform configuration starts by initializing the Akamai provider and by providing the authentication credentials. After that, the configuration uses the [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list) resource and the following block to create the network list:

```
resource "akamai_networklist_network_list" "network_list" {
  name        = "Documentation Network"
  type        = "IP"
  description = "This is a test network used by the documentation team."
  list        = ["192.168.1.1","192.168.1.2","192.168.1.3"]
  mode        = "REPLACE"
}
```

The properties used in this block are described in the following table:

| Argument    | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| name        | Name of the new network list. Names don't have to be unique: you can have multiple network lists that share the same name. However, when the list is created it will be issued a unique ID, a value comprised of a numeric prefix and the list name. (Or a variation of that name. For example, a list named **Documentation Network** will be given an ID similar to **108970_DOCUMENTATIONNETWORK**, with the blank space in the name being removed. |
| type        | Indicates the type of addresses used on the list. Allowed values are:<br /><br />* **IP**. For IP/CIDR addresses.<br />* **GEO**. For ISO 3166 geographic codes.<br /><br />Note that you cannot mix IP/CIDR addresses and geographic codes on the same list. |
| description | Brief description of the network list.                       |
| list        | Array containing either the IP/CIDR addresses or the geographic codes to be added to the new network list. For example:<br /><br />`list = ["US", "CA", "MX"]`<br /><br />Note that the list value must be formatted as an array even if you are only adding a single item to that list:<br /><br />`list = ["US"]`<br /><br />Note, too that `list` is the one optional argument available to you: you don't have to include this argument in your configuration. However, leaving out the `list` argument also means that you'll be creating a network list that has no IP/CIDR addresses or geographic codes. |
| mode        | Specifies whether the addresses/geographic codes on the list should:<br /><br />* Be added to the existing set of addresses (**APPEND**).<br />* Be removed from the existing set of addresses (**REMOVE**).<br />* Replace the existing set of addresses (**REPLACE**).<br /><br />Because we're creating a new network list, there isn't going to be an existing set of addresses. Consequently, you can set the mode either to **APPEND** or to **REPLACE**.<br /> |

After you've created your configuration, run `terraform plan` to help validate and verify that configuration. Assuming that command succeeds, you can then run `terraform apply` to create the new network list. If the apply command is successful, you'll see output similar to this; that output includes the ID (**108972_DOCUMENTATIONNETWORK** ) of the new network list:

```
akamai_networklist_network_list.network_list: Creation complete after 4s [id=108972_DOCUMENTATIONNETWORK]<br />Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```



## Activating a Network List

After the network list has been created, we then use this Terraform block to activate the network:

```
resource "akamai_networklist_activations" "activation" {
  network_list_id     = akamai_networklist_network_list.network_list.uniqueid
  network             = "STAGING"
  notes               = "This is the staged network for the documentation test network."
  notification_emails = ["gstemp@akamai.com"]
}
```

Note that this step is optional: you don't have to activate a network immediately after creating the network. Instead, you can call the [akamai_networklist_activations](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_activations) resource at any time to activate a network. Note that a network can't be referenced by Akamai products or by other Terraform resources until it's been activated.

Our activation block uses the following arguments:

| Argument            | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| network_list_id     | Unique identifier of the network being activated. In this example, **akamai_networklist_network_list.network_list.uniqueid** represents the ID of the network that we just created. To activate a previously-existing network, specify the ID of that network. |
| network             | Indicates which network is to be activated. Allowed bvalues are:<br /><br />* STAGING<br />* PRODUCTION<br /><br />Note that this argument is optional. If not included in your configuration, the staging network will automatically be selected for activation. |
| notes               | Optional description of the network and/or information about the activation. |
| notification_emails | JSON array containing the email addresses of users who will be notified when the network activation finishes. |


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

In addition, each user included in the notification_emails argument will be emailed an activation notice.

## Modifying an Existing Network List

Modifying an existing network list can be tricky, if only because you need to know which Terraform resource is used to modify which property values. Use the following table to help you determine which resource is the appropriate one for you to use:

| Property    | Recommended Resource for Modifying the Property              |
| ----------- | ------------------------------------------------------------ |
| name        | [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list)<br /><br />By default, the network name is included in the network ID; for example, a network named **USANetwork** will have an ID similar to **110987_USANETWORK**. However, if you change the name of the network the network ID will not change. For example, you can change the name **USANetwork** to **HomeNetwork**, but the ID for that network will still be **110987_USANETWORK**. |
| type        | [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list)<br /><br />If you change the `type` I(e.g., if you switch from an IP list to a GEO list) you will need to change the list value as well. (Unless, of course, the list is currently empty.) |
| description | [akamai_networklist_description](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_description)<br /><br />Note that you can also use the [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list) resource to change a network list description. |
| list        | [akamai_networklist_network_list](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/networklist_network_list)<br /><br />Be sure and set the `mode` to the appropriate operation type (**APPEND**, **REPLACE**, **REMOVE**) when changing the list. |
| recipients  | akamai_networklist_subscription<br /><br />To have email notifications sent to multiple addresses, separate those addresses by using commas:<br /><br />[`"gstemp@akamai.com", "karim.nafir@mail.com"]` |

For example, this configuration changes the description assigned to the network list named **Documentation Network**:

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

data "akamai_networklist_network_lists" "network_lists_filter" {
  name = "Documentation Network"
}

resource "akamai_networklist_description" "network_list_description" {
  network_list_id = data.akamai_networklist_network_lists.network_lists_filter.uniqueid
  name            = "Documentation Network"
  description     = "Updated description for the documentation network."
}
```

And this configuration changes the set of email addresses assigned to that same. network list:

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

data "akamai_networklist_network_lists" "network_lists_filter" {
  name = "Documentation Network"
}

resource "akamai_networklist_subscription" "subscribe" {
  network_list = data.akamai_networklist_network_lists.network_lists_filter.list
  recipients   = ["gstemp@akamai.com", "gmstemp@hotmail.com"]
}
```

## Deactivating a Network List

To deactivate a network list, just remove all the addresses/geographic codes that have been assigned to that list. For example, this Terraform configuration deactivates the network list **Documentation Network**:

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

These two lines remove the addresses/geographic codes from the list:

```
list = []
mode = "REPLACE"
```

The first line simply sets the list (i.e., the collection of device addresses or geographic codes) to an empty array. The second line then specifies that the new, empty list should replace the existing set of addresses/geographic codes. In effect, that removes all the IP addresses/geographic codes from the network list, which, in turn, deactivates that list.
