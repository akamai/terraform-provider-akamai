provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_properties" "akaproperties" {
  group = "test"
  contract = "ctr_test"
}