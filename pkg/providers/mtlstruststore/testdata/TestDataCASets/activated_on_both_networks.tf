provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_sets" "test" {
  activated_on = "staging+production"
}