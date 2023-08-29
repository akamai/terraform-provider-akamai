resource "akamai_clientlist_list" "test_list" {
  name        = "List Name Updated"
  tags        = ["a", "c", "d"]
  notes       = "List Notes Updated"
  type        = "ASN"
  contract_id = "12_ABC"
  group_id    = 12

  items {
    value       = "1"
    description = "Item 1 Desc"
    tags        = ["item1Tag1", "item1Tag2"]
  }
  items {
    value       = "12"
    description = "Item 12 Desc Updated"
    tags        = ["item12Tag1", "item12Tag2"]
  }
  items {
    value       = "1234"
    description = "Item 1234 Desc"
    tags        = ["1234Tag"]
  }
}
