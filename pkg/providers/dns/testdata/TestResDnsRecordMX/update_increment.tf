provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "record" {
  zone               = "exampleterraform.io"
  name               = "exampleterraform.io"
  recordtype         = "MX"
  ttl                = 300
  priority           = 3
  priority_increment = 3
  target = [
    "mx1.test.com.",
    "mx2.test.com.",
    "mx3.test.com.",
  ]
}

