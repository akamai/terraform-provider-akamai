provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id            = 500051
  staging_generation_id = 6300
  # re-adding the same value production was already tracked at is recognized as a no-op: no API call is made.
  production_generation_id = 6200
}
