provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_contracts" "akacontracts" {
}

output "aka_contract_id1" {
  value = data.akamai_contracts.akacontracts.contracts[0].contract_id
}

output "aka_contract_id2" {
  value = data.akamai_contracts.akacontracts.contracts[1].contract_id
}

output "aka_contract_typ_name1" {
  value = data.akamai_contracts.akacontracts.contracts[0].contract_type_name
}

output "aka_contract_typ_name2" {
  value = data.akamai_contracts.akacontracts.contracts[1].contract_type_name
}
