provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name        = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  contract_id = "ctr_0"
  group_id    = "grp_0"
}
