provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_activation" "test" {
  lineage_id               = 500042
  production_generation_id = 8100

  timeouts {
    create = "5ms"
  }
}
