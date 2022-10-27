provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_datastreams" "test" {
  group_id = "asdf_1234"
}