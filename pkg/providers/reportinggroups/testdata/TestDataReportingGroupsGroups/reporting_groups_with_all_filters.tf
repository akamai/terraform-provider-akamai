provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_groups" "test" {
  contract_id          = "ctr_123"
  group_id             = 456
  reporting_group_name = "First Reporting Group"
  cp_code_id           = "111"
}

