provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_hostname_bucket" "test" {
  property_id = "prp_111"
  network     = "STAGING"
}
