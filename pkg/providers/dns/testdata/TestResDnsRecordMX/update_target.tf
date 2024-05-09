provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "MX"
  ttl        = 300
  target = [
    "5 mx1.test.com.",
    "10 mx2.test.com.",
    "15 mx3.test.com.",
    "20 mx4.test.com.",
  ]
}

