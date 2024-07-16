provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "aaaa_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "AAAA"
  ttl        = 300
  target     = ["18.244.102.124"]
}

