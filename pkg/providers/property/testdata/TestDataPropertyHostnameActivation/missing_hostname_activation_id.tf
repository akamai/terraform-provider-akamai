provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostname_activation" "activation" {
  property_id = "1"
}