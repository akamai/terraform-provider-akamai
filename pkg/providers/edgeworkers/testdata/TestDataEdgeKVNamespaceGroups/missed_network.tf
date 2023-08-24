provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgekv_groups" "test" {
  namespace_name = "test_namespace"
}