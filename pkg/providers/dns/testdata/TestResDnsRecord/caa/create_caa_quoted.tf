provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "caa_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "CAA"
  ttl        = 300
  target     = ["0 issue \"ca.example.net\""]
}
