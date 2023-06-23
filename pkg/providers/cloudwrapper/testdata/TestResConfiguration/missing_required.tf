provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudwrapper_configuration" "test" {
  location {
    capacity {

    }
  }
}