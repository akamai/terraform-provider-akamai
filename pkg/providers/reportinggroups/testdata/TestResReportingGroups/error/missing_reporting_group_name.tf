provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_reportinggroups_group" "test" {
  access_group = {
    contract_id = "test_contract"
    group_id    = "12345"
  }
  contract = {
    contract_id = "test_contract"
    cp_codes = [
      {
        cp_code_id = "111111"
      }
    ]
  }
}