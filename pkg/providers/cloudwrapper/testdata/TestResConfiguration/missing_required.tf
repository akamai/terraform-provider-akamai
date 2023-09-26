provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudwrapper_configuration" "test" {
  location {
    capacity {

    }
  }
}