data "akamai_edgekv_group_items" "test" {
  namespace_name = "test_namespace"
  network        = "incorrect_network"
  group_name     = "TestGroup"
}