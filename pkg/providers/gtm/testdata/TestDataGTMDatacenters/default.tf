provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_datacenters" "test" {
  domain = "test.domain.com"
}
