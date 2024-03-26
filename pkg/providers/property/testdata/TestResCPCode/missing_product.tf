provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cp_code" "test" {
  name        = "test cpcode"
  group_id    = "grp_2"
  contract_id = "ctr_1"
}
