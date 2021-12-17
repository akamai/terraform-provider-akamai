provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_networklist_network_lists" "test" {
  uniqueid = "86093_AGEOLIST"
}
