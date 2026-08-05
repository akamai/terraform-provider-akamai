provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id               = 500091
  staging_generation_id    = 9010
  production_generation_id = 9021
}
