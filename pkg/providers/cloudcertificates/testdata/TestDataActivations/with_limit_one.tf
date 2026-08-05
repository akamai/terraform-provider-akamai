provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_activations" "test" {
  lineage_id = 500022
  limit      = 1
}
