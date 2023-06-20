provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract_id   = "2"
  group_id      = "2"
  product_id    = "2"
  edge_hostname = "test2.edgesuite.net"
  certificate   = 123
  ip_behavior   = "IPV6_COMPLIANCE"
}

output "edge_hostname" {
  value = akamai_edge_hostname.edgehostname.edge_hostname
}