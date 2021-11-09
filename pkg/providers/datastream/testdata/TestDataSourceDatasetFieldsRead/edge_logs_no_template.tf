provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_datastream_dataset_fields" "test" {}