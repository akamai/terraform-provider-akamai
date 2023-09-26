provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_clientlist_list" "test_list" {
  name        = "List Name"
  tags        = ["a", "b"]
  notes       = "List Notes"
  type        = "ASN"
  contract_id = "12_ABC"
  group_id    = 12

  items {
    value       = "12"
    description = "Item 12 Desc"
    tags        = ["item12Tag1", "item12Tag2"]
  }
  items {
    value = "12"
  }
}
