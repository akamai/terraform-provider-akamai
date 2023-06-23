provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_clientlist_list" "test_list" {
  name        = "List Name Updated"
  tags        = ["a", "c", "d"]
  notes       = "List Notes Updated"
  type        = "ASN"
  contract_id = "12_ABC"
  group_id    = 12
}
