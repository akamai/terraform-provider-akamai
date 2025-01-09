provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_datacenter" "test" {
  domain = "test_domain"
}
