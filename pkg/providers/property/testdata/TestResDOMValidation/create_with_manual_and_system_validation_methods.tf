provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_validation" "test" {
  domains = [
    {
      domain_name       = "example1.com"
      validation_scope  = "HOST"
      validation_method = "SYSTEM"
    },
    {
      domain_name       = "example2.com"
      validation_scope  = "DOMAIN"
      validation_method = "MANUAL"
    },
  ]
}
