provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudwrapper_properties" "test" {
  unused = true
}