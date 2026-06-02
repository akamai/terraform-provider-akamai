provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_clientlist_list" "test_list" {
  name        = "List Name"
  tags        = ["a", "b"]
  notes       = "List Notes"
  type        = "REQUEST_HEADER_NAME_VALUE"
  contract_id = "12_ABC"
  group_id    = 12

  items {
    key    = "header1"
    values = ["val1"]
  }
  items {
    key    = "header1"
    values = ["val2"]
  }
}

