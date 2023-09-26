provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_edgekv_groups" "test" {
  network = "staging"
}