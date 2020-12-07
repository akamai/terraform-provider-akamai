provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract = "ctr_2"
  group = "grp_2"
  product = "prd_2"
  edge_hostname = "test.akamaized.net"
  ip_behavior = "unknown"
}
