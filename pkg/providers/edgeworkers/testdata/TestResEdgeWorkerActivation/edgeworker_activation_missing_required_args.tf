provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edgeworker_activation" "test" {
}