provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edge_hostnames" "test" {
  group_id = "grp_123"
}