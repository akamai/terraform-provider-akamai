provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_activation_status" "test" {
  lineage_id    = 500019
  activation_id = 100
}
