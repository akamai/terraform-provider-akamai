provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_group" "akagroup" {
  group_name = "Example.com-1-1TJZH5"
  contract_id = "ctr_1234"
}