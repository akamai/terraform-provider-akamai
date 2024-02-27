provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edgeworkers_activation" "test" {
  edgeworker_id = 1234
  network       = "PRODUCTION"
  version       = "test1"
  note          = "note for edgeworkers activation updated"
}