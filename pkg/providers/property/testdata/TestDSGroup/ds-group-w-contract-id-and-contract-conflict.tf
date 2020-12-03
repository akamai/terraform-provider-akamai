provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_group" "akagroup" {
  name = "Example.com-1-1TJZH5"
  contract = "ctr_1234"
  contract_id = "ctr_1234"
}