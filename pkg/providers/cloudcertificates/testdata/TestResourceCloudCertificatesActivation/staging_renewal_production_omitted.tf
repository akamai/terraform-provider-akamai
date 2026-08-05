provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id            = 500051
  staging_generation_id = 6300
  # production_generation_id intentionally omitted: this is a staging-only renewal, freezing production at
  # whatever it was already tracked at.
}
