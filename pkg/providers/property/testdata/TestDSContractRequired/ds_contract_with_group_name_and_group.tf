provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_contract" "akacontract" {
  group_id   = "grp_12345"
  group_name = "Example.com-1-1TJZH5"
}

output "aka_contract" {
  value = data.akamai_contract.akacontract.id
}