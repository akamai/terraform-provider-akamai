provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edge_hostname" "importedgehostname" {
  contract_id   = "ctr_1"
  group_id      = "grp_2"
  edge_hostname = "test.akamaized.net"
  ip_behavior   = "IPV4"
}
