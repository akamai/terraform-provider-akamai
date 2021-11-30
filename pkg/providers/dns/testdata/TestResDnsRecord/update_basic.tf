provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_dns_record" "a_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "A"
  active     = true
  ttl        = 300
  target     = ["10.0.0.4", "10.0.0.5"]
}
