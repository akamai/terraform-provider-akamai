provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlstruststore_ca_set_activation" "test" {
  version = 1
  network = "STAGING"
}