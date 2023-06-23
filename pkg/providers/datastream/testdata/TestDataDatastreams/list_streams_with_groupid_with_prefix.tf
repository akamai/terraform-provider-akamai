provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_datastreams" "test" {
  group_id = 1234
}