provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edge_hostname" "createhostname" {
  contract      = "ctr_1"
  group         = "grp_2"
  product       = "prd_2"
  edge_hostname = "test.akamaized.net"
  ip_behavior   = "IPV4"
}
