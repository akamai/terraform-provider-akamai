provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_datacenters" "test" {
  domain = "test.domain.com"
}
