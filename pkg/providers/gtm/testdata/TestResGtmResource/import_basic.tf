provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_resource" "test" {
  aggregation_type = "latest"
  domain           = "test_domain"
  name             = "tfexample_resource_1"
  type             = "XML load object via HTTP"
}
