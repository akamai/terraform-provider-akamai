provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_datastream_dataset_fields" "test" {
  product_id = "PROD_1"
}