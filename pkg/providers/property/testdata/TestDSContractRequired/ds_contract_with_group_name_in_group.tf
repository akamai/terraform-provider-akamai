provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_contract" "akacontract" {
  group = "Example.com-1-1TJZH5"
}

output "aka_contract" {
  value = data.akamai_contract.akacontract.id
}