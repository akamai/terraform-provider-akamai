provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {
  id   = "12345"
  name = "test_name"
}