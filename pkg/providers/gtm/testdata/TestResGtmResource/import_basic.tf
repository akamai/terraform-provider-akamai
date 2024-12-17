provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_resource" "test" {
  aggregation_type = "latest"
  domain           = "gtm_terra_testdomain.akadns.net"
  name             = "tfexample_resource_1"
  type             = "XML load object via HTTP"
}
