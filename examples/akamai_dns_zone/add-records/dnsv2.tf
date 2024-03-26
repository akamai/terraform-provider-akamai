terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
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
