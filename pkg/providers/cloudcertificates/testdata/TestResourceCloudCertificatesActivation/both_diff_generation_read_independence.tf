provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id               = 500093
  staging_generation_id    = 9800
  production_generation_id = 9801
}
