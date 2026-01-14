provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_clientlist_list" "test_list" {
  name        = "List Name"
  tags        = ["a", "b"]
  notes       = "List Notes"
  type        = "DOMAIN"
  contract_id = "12_ABC"
  group_id    = 12

  items {
    value       = "a.com"
    description = "Domain a"
  }
  items {
    value       = "b.com"
    description = "Domain b"
  }
  items {
    value       = "c.com"
    description = "Domain c"
  }
}

output "version" {
  value = akamai_clientlist_list.test_list.version
}
