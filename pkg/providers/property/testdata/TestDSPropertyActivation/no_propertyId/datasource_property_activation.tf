provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_activation" "test" {
  version = 1
}