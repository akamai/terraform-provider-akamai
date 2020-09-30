terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {}

resource "akamai_dns_record" "a_record" {
  zone       = "akavdev.net"
  name       = "akavdev.net"
  recordtype = "A"
  active     = true
  ttl        = 300
  target     = ["10.0.0.2", "10.0.0.3"]
}
