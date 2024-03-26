provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_edgekv_group_items" "test" {
  namespace_name = "test_namespace"
  network        = "staging"
  group_name     = "1234"
  items = {
    key2 = "updatedValue"
    key3 = "value3"
    key4 = "value4"
  }
}

