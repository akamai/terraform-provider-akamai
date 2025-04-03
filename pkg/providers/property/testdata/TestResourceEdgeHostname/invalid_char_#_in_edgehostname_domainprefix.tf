provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract_id   = "ctr_2"
  group_id      = "grp_2"
  product_id    = "prd_2"
  edge_hostname = "tes#t.akamaized.net"
  ip_behavior   = "IPV6_PERFORMANCE"
}
