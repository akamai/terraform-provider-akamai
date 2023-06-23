provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edge_hostname" "createhostname" {
  contract_id   = "ctr_1"
  group_id      = "grp_2"
  product_id    = "prd_2"
  edge_hostname = "test.akamaized.net"
  ip_behavior   = "IPV4"
}
