provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_networklist_network_lists" "test" {
  network_list_id = "86093_AGEOLIST"
}
