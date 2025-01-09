provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_domain" "test" {
  name = "test_domain"
  type = "weighted"
}
