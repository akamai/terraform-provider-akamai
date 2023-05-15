provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_datastreams" "test" {
  group_id = 1234
}