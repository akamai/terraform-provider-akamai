provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_datastream_activation_history" "test" {
  stream_id = 7050
}