provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "aaaa_record_full" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "AAAA"
  ttl        = 300
  target     = ["1000:0000:0000:0000:0000:0000:0000:0001", "1000:0000:0000:0000:0000:0000:0000:0002", "1000:0000:0000:0000:0000:0000:0000:0003", "1000:0000:0000:0000:0000:0000:0000:0004"]
}
