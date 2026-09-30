provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id               = 500060
  staging_generation_id    = 7000
  production_generation_id = 7001
}
