resource "akamai_clientlist_list" "test_list" {
  name        = "List Name"
  tags        = []
  notes       = "List Notes"
  type        = "IP"
  contract_id = "12_ABC"
  group_id    = 12
}
