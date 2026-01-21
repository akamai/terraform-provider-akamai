provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}
data "akamai_property_account_hostnames" "test" {
  contract_id = "ctr_1"
  group_id    = "grp_1"
  hostname    = "example.com"
  cname_to    = "example.com.edgesuite.net"
  network     = "PRODUCTION"
  sort        = "hostname:a"
}