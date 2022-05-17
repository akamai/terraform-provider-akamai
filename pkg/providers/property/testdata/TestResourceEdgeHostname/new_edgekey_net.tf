provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract      = "ctr_2"
  group         = "grp_2"
  product       = "prd_2"
  edge_hostname = "test.edgekey.net"
  certificate   = 123
  ip_behavior   = "IPV6_PERFORMANCE"
}

output "edge_hostname" {
  value = akamai_edge_hostname.edgehostname.edge_hostname
}