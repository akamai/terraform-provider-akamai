provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_domains" "test" {
  domains = [
    {
      domain_name      = "test1.com"
      validation_scope = "HOST"
    },
    {
      domain_name      = "test1.com"
      validation_scope = "DOMAIN"
    },
  ]
}
