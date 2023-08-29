resource "akamai_clientlist_list" "test_list" {
  name        = "List Name"
  tags        = ["a", "b"]
  notes       = "List Notes"
  type        = "ASN"
  contract_id = "12_ABC"
  group_id    = 12

  items {
    value       = "1"
    description = "Item 1 Desc"
    tags        = ["item1Tag1", "item1Tag2"]
  }
  items {
    value           = "123"
    expiration_date = "2026-12-26T01:00:00+00:00"
  }
  items {
    value       = "12"
    description = "Item 12 Desc"
    tags        = ["item12Tag1", "item12Tag2"]
  }
}

output "version" {
  value = akamai_clientlist_list.test_list.version
}
