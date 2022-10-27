provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_datastreams" "test" {}