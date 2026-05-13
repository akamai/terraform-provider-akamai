provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_reportinggroups_group" "test" {
  reporting_group_name = "test-reporting-group"
  access_group = {
    contract_id = "ctr_test_contract"
    group_id    = "grp_12345"
  }
  contract = {
    contract_id = "ctr_test_contract_2"
    cp_codes = [
      {
        cp_code_id = "cpc_111111"
      }
    ]
  }

}