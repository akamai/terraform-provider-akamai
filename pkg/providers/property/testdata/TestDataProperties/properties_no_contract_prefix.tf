provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_properties" "akaproperties" {
  group = "grp_test"
  contract = "test"
}