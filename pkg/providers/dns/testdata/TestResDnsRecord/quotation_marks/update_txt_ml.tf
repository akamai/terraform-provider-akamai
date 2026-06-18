provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_record" "txt_record" {
  zone       = "infrastructure.domain.net"
  name       = "infrastructure.domain.net"
  recordtype = "TXT"
  target     = ["v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all", "\"foo\" \"foo\""]
  ttl        = 1800
}
