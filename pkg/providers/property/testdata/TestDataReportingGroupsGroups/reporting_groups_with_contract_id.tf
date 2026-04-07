provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_reportinggroups_groups" "test" {
  contract_id = "ctr_456"
}

