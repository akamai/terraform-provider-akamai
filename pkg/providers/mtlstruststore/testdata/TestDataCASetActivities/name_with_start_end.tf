provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {
  name  = "test name"
  start = "2024-04-16T12:08:34.099457Z"
  end   = "2025-04-16T12:08:34.099457Z"
}