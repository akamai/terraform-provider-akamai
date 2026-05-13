provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_cp_codes" "test" {
  contract_id = "1-2ABCDE"
}

