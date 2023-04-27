provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edgekv_group_items" "test" {
  namespace_name = "test_namespace"
  network        = "staging"
  group_name     = "1234"
  items = {
    key2 = "value2"
    key3 = "value3"
    key1 = "value1"
  }
}

