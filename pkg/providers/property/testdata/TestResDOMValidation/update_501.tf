provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_validation" "test" {
  domains = [
    for i in range(0, 501) : {
      domain_name      = "update-test${i}.example.com"
      validation_scope = "HOST"
    }
  ]
}
