provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_zone_dnssec_status" "test" {
  zone = "test.zone.net"
}
