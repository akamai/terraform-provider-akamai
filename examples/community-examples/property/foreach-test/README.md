# An example for SaaS Providers

This is an example of how a SaaS service provider would use a single configuration template for many similar configs by using a new Terraform feature called "for_each". The basic function of "for_each" is to allow you to iterate over an array or a map and provision infrastructure for each of them. This is ideal for our use-case. 

In our example, our SaaS service requires that each instance differs by username/password and hostname. To drive this, we can define a complex object type variable that describes each of our customers SaaS instances like so:-

```
variable "customers" {
        type = map(object({
                username = string
                password = string
        }))
}
```

and then we can populate this variable with our terraform.tfvars

```
customers =  {
     "foreach1.wheep.co.uk" = {
                username = "test"
                password = "test"
        },
     "foreach2.wheep.co.uk" = {
                username = "test2"
                password = "test2"
        }
}
```

We then need to modify our main.tf to include the "for_each" logic in each resource that needs to be individually provisioned for each instance.

For example

```
data "template_file" "rules" {
        for_each = var.customers

        template = data.template_file.rule_template.rendered
        vars = {
                username = each.value.username
                password = each.value.password
        }
}
```

Given our configuration for the "customers" variable, this will cause Terraform to create 2 instances of "template_file.rules". Each will be referenced with the key value of our variable. One would be "template_file.rules[foreach1.wheep.co.uk]" and the other "template_file.rules[foreach2.wheep.co.uk]". We don't need to concern ourselves too much what the key value is, only that the key value exists and is required when referencing it. In practical terms, this means that you need to supply the key value when you reference this resource from another resource.

```
resource "akamai_property" "property" {

  for_each = var.customers

  name        = each.key
  cp_code     = akamai_cp_code.cpcode[each.key].id
  contract = data.akamai_contract.contract.id
  group = data.akamai_group.group.id
  product     = "prd_Site_Accel"
  rule_format = "v2018-02-27"

  hostnames    = {
        "${each.key}" = akamai_edge_hostname.edge_hostname.edge_hostname
  }
  rules       = data.template_file.rules[each.key].rendered
  is_secure = true

}
```

You can see this in action when you do a "terraform apply".

```
akamai_property_activation.activation["foreach1.wheep.co.uk"]: Modifying... [id=atv_8064436]
akamai_property_activation.activation["foreach2.wheep.co.uk"]: Modifying... [id=atv_8064438]
akamai_property_activation.activation["foreach1.wheep.co.uk"]: Modifications complete after 5s [id=atv_8064436]
akamai_property_activation.activation["foreach2.wheep.co.uk"]: Still modifying... [id=atv_8064438, 10s elapsed]
...
akamai_property_activation.activation["foreach2.wheep.co.uk"]: Modifications complete after 2m19s [id=atv_8064592]
```

You can also reference a particular instance in case you want to target something specifically. For example...

```
terraform destroy -target='akamai_property_activation.activation["foreach1.wheep.co.uk"]'
...

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # akamai_property_activation.activation["foreach1.wheep.co.uk"] will be destroyed
  - resource "akamai_property_activation" "activation" {
      - activate = true -> null
      - contact  = [
          - "icass@akamai.com",
        ] -> null
      - id       = "atv_8064436" -> null
      - network  = "STAGING" -> null
      - property = "prp_594500-0aa0a95548921426a2d416e373c7354a48218ffa" -> null
      - status   = "ACTIVE" -> null
      - version  = 1 -> null
    }
```

