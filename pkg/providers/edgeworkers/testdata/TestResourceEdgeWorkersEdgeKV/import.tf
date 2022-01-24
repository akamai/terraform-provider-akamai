provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_edgekv" "test" {
}