provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract = "ctr_2"
  group = "grp_2"
  product = "prd_2"
  edge_hostname = "test.edgesuite.net"
  certificate = 123
  ipv4 = true
  ipv6 = true
}

