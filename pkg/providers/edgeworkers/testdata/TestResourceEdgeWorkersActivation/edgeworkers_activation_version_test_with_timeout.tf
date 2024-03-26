provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edgeworkers_activation" "test" {
  edgeworker_id = 1234
  network       = "STAGING"
  version       = "test"
  note          = "note for edgeworkers activation"
  timeouts {
    default = "2h"
    delete  = "3h"
  }
}