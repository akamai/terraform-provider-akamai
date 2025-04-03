provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostname_activation" "activation" {
  hostname_activation_id = "1"
}