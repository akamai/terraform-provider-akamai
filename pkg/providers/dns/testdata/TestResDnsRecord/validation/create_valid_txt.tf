provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "txt_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "TXT"
  ttl        = 300
  target     = ["valid-target"]
}