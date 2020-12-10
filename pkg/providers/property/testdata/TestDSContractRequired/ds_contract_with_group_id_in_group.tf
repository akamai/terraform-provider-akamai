provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_contract" "akacontract" {
  group = "grp_12345"
}

output "aka_contract" {
  value = data.akamai_contract.akacontract.id
}