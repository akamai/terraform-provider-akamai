provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_clientlist_lists" "lists" {
  name = "test"
  type = ["REQUEST_HEADER_NAME_VALUE"]
}

