provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edge_hostname" "test" {
  id          = "ehn_123"
  contract_id = "ctr_123"
  group_id    = "grp_123"
}
