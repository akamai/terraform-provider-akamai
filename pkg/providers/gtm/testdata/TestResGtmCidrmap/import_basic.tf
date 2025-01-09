provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_cidrmap" "test" {
  domain = "test_domain"
  name   = "tfexample_cidrmap_1"
}