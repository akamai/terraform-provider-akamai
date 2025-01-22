provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "srv_record" {
  zone       = "origin.org"
  name       = "origin.example.org"
  recordtype = "SRV"
  ttl        = 300
  target = [
    "10 60 5060 big.example.com",
    "10 40 5060 small.example.com",
    "20 100 5060 tiny.example.com"
  ]
}
