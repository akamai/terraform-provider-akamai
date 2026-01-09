provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostname_audit_history" "test" {
  hostname = "example.com"
}
