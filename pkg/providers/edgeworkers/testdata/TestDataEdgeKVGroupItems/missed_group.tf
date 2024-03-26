provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgekv_group_items" "test" {
  namespace_name = "test_namespace"
  network        = "staging"
}