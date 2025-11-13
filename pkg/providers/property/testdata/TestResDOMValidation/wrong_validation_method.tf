provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_validation" "test" {
  domains = [
    {
      domain_name       = "test1.example.com"
      validation_scope  = "HOST"
      validation_method = "WRONG"
    }
  ]
}
