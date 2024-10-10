provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract_id   = "ctr_2"
  group_id      = "grp_2"
  product_id    = "prd_2"
  edge_hostname = "test.edgekey.net"
  certificate   = 800800
  ip_behavior   = "IPV6_PERFORMANCE"
}

output "edge_hostname" {
  value = akamai_edge_hostname.edgehostname.edge_hostname
}