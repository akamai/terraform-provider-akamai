provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "akamaicdn_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "AKAMAICDN"
  ttl        = 300
  target     = ["xyz-test.edgesuite.net"]
}
