provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_group" "akagroup" {
  name        = "Example.com-1-1TJZH5"
  contract_id = "ctr_1234"
}