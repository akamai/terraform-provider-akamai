provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_contract" "akacontract" {
  group = "Example.com-1-1TJZH5"
  group_name = "Example.com-1-1TJZH5"
}

output "aka_contract" {
  value = data.akamai_contract.akacontract.id
}