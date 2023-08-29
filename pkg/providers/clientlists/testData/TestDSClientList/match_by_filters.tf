provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_clientlist_lists" "lists" {
  name = "test"
  type = ["IP", "GEO"]
}