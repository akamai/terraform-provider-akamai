provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_datastream_activation_history" "test" {
  stream_id = 7051
}