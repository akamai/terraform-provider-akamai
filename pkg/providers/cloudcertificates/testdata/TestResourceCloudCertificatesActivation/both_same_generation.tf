provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id               = 500050
  staging_generation_id    = 6000
  production_generation_id = 6000
}
