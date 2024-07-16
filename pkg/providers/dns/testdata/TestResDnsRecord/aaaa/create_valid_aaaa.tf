provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "aaaa_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "AAAA"
  ttl        = 300
  target     = ["2001:db8::68", "::ffff:192.0.2.1"]
}

