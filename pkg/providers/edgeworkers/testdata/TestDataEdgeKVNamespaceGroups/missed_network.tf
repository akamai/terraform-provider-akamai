provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_edgekv_groups" "test" {
  namespace_name = "test_namespace"
}