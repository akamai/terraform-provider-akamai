provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_group" "akagroup" {
  name = "group-example.com"
  contract = "ctr_1234"
}