provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_group" "akagroup" {
  name = "group-example.com"
  contract_id = "ctr_1234"
}