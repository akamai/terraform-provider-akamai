provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_imaging_policy_video" "policy" {
  policy_id    = "test_policy"
  contract_id  = "test_contract"
  policyset_id = "test_policy_set"
  policy {
  }
}