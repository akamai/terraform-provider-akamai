provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_clientlist_list" "test_list" {
  name        = "List Name"
  tags        = ["a", "b"]
  notes       = "List Notes"
  type        = "USER_ID"
  contract_id = "12_ABC"
  group_id    = 12

  items {
    value       = "3a453537-faa8-4525-b5db-022447bbbf2a"
    description = "Item 1 Desc"
    tags        = ["item1Tag1", "item1Tag2"]
  }
  items {
    value           = "07e29045-7739-4bd9-8cfb-9f118e000337"
    expiration_date = "2026-12-26T01:00:00+00:00"
  }
  items {
    value       = "e164394a-5ae1-4208-8487-1ac0f368ecf3"
    description = "Item 3 Desc"
    tags        = ["item3Tag1", "item3Tag2"]
  }
}

output "version" {
  value = akamai_clientlist_list.test_list.version
}
