provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edge_hostnames" "test" {
  contract_id = "ctr_123"
}