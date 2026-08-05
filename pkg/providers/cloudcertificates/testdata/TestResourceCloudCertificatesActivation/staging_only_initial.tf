provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id            = 500092
  staging_generation_id = 9030
}
