provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_reportinggroups_group" "test" {
  reporting_group_name = "test-reporting-group"
  access_group = {
    contract_id = "test_contract"
    group_id    = "54321"
  }
  contract = {
    contract_id = "test_contract_2"
    cp_codes = [
      {
        cp_code_id = "111111"
      }
    ]
  }
}