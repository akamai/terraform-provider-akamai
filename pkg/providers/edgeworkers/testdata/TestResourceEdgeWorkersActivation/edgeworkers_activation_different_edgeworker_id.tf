provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edgeworkers_activation" "test" {
  edgeworker_id = 4321
  network       = "STAGING"
  version       = "test"
  note          = "note for edgeworkers activation"
}