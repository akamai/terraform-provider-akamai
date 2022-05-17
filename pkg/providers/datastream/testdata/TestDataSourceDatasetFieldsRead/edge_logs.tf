provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_datastream_dataset_fields" "test" {
  template_name = "EDGE_LOGS"
}