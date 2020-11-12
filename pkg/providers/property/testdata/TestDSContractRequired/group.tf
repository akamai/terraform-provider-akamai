provider "akamai" {
  edgerc = "~/.edgerc"
  property_section = "afadencloud-example-01"
}

data "akamai_contract" "akacontract" {
  group = var.group
}

data "akamai_group" "akagroup" {
  name = var.group
  contract_id = data.akamai_contract.akacontract.id
}

data "akamai_group" "akgroup" {
  name = var.group
  contract = data.akamai_contract.akacontract.id
}

variable "group" {
  description = "Name of the group associated with this property"
  type = string
  default = "test"
}

output "aka_contract" {
  value = data.akamai_contract.akacontract.id
}

output "aka_group" {
  value = data.akamai_group.akagroup.id
}
output "ak_group" {
  value = data.akamai_group.akgroup.id
}


