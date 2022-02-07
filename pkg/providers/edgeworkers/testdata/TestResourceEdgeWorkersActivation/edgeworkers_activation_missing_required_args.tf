provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edgeworkers_activation" "test" {
}