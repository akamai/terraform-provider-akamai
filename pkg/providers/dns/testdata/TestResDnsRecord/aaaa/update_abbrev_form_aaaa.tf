provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "aaaa_record_abbrev" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "AAAA"
  ttl        = 300
  target     = ["1000:0:0:0:0:0:0:1", "1000:0:0:0:0:0:0:2", "1000:0:0:0:0:0:0:3", "1000:0:0:0:0:0:0:4"]
}
