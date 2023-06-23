provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_datastream_activation_history" "test" {
  stream_id = 7050
}