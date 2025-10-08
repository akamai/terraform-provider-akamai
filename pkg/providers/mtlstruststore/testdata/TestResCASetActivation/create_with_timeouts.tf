provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlstruststore_ca_set_activation" "test" {
  ca_set_id = "12345"
  version   = 1
  network   = "STAGING"
  timeouts = {
    create = "200ms"
    update = "200ms"
    delete = "100ms"
  }
}