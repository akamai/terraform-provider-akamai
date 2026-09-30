provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id               = 500051
  staging_generation_id    = 6100
  production_generation_id = 6200
}
