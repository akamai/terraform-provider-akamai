provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudwrapper_locations" "test" {}