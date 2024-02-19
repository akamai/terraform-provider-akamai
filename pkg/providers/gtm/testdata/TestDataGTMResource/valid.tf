provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_resource" "my_gtm_resource" {
  domain        = "test.domain.net"
  resource_name = "resource1"
}