provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract = "2"
  group = "2"
  product = "2"
  edge_hostname = "test2.edgesuite.net"
  certificate = 123
  ip_behavior = "IPV6_COMPLIANCE"
}

