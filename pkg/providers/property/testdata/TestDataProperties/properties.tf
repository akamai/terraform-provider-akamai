provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_properties" "akaproperties" {
  group_id    = "grp_test"
  contract_id = "ctr_test"
}