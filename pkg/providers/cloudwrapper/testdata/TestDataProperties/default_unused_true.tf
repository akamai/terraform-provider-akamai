provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudwrapper_properties" "test" {
  unused = true
}