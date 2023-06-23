provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgekv_group_items" "test" {
  network    = "staging"
  group_name = "TestGroup"
}