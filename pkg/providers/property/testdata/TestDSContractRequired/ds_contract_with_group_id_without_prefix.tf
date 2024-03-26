provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_contract" "akacontract" {
  group_id = "12345"
}

output "aka_contract" {
  value = data.akamai_contract.akacontract.id
}