provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edge_hostname" "edgehostname" {
  contract = "ctr_2"
  group = "grp_2"
  product = "prd_2"
  edge_hostname = "test.aka"
  ip_behavior = "IPV4"
}

output "edge_hostname" {
  value = akamai_edge_hostname.edgehostname.edge_hostname
}
