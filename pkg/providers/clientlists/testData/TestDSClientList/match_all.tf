provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_clientlist_lists" "lists" {}

