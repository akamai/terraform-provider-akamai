provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_dns_zone" "test_without_group" {
  contract       = "ctr1"
  zone           = "primaryexampleterraform.io"
  type           = "primary"
  comment        = "This is a test primary zone"
  sign_and_serve = false
}