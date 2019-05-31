provider "akamai" {
    edgerc = "/root/.edgerc"
    dns_section = "dns"
}

locals {
  zone = "akavdev.net"
}

resource "akamai_dns_zone" "test_zone" {
    contract = "C-1FRYVV3"
    zone = "${local.zone}"
    #type = "SECONDARY"
    masters = ["1.2.3.4" , "1.2.3.5"]
    type = "PRIMARY"
    comment =  "This is a test zone"
    group     = "64867"
    signandserve = false
}


data "akamai_authorities_set" "ns" {
  contract = "C-1FRYVV3"
}


output "authorities" {
  value = "${join(",", data.akamai_authorities_set.ns.authorities)}"
}


