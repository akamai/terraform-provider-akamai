provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "aaaa_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "AAAA"
  ttl        = 300
  target     = ["1111:2222:3333:4444:55555:6666:7777:8888"]
}

