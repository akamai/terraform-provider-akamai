provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudwrapper_configurations" "test" {}