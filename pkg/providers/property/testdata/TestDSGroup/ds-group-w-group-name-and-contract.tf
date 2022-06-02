provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_group" "akagroup" {
  group_name = "Example.com-1-1TJZH5"
  contract   = "ctr_1234"
}