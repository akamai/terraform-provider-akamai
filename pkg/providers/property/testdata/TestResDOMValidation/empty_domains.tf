provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_validation" "test" {
  domains = []
}
