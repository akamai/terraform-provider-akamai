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
    key         = "header1"
    values      = ["val1"]
    description = "Header 1 Desc Updated"
  }
  items {
    key         = "header2"
    values      = ["val2"]
    description = "Header 2 Desc"
  }
  items {
    key         = "header4"
    values      = ["val4"]
    description = "Header 4 Desc"
  }
}
