provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_bindings" "test" {
  lineage_id = 500005
  network    = "NOT_A_NETWORK"
}
