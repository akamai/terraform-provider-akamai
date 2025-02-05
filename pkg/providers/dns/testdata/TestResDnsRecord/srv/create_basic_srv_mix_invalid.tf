provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "srv_record" {
  zone       = "origin.org"
  name       = "origin.example.org"
  recordtype = "SRV"
  ttl        = 300
  priority   = 10
  weight     = 60
  port       = 5060
  target = [
    "big.example.com",
    "small.example.com",
    "20 50 5060 tiny.example.com"
  ]
}
