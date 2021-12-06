provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_dns_zone" "test_secondary_zone" {
  contract       = "ctr1"
  zone           = "secondaryexampleterraform.io"
  masters        = ["1.2.3.4", "1.2.3.5"]
  type           = "secondary"
  comment        = "This is a secondary test zone"
  sign_and_serve = false
}