provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudwrapper_configuration" "test" {}