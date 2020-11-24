provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_group" "akagroup" {
  group_name = "group-example.com"
  contract = "ctr_1234"
}