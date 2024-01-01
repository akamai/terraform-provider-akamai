provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_resources" "my_gtm_resources" {
  domain        = "test.domain.net"
}