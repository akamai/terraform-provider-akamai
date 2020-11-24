provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_group" "akagroup" {
  group_name = "group-example.com"
  contract_id = "ctr_1234"
}