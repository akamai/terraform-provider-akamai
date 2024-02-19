provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_asmap" "my_gtm_asmap" {
  domain = "test.domain.net"
}
