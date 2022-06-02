provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_contract" "akacontract" {
  group = "12345"
}

output "aka_contract" {
  value = data.akamai_contract.akacontract.id
}