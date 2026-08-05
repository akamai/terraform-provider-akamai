provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id               = 500098
  staging_generation_id    = 7200
  production_generation_id = 7300
}
