provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cp_code" "test" {
  name        = "test cpcode"
  contract_id = "ctr1"
  group       = "grp1"
  product     = "prd2"
}
