provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edgekv_group_items" "test" {
  namespace_name = "test_namespace"
  network        = "staging"
  items = {
    key1 = "value1"
    key2 = "value2"
  }
}

