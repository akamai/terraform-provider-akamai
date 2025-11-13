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
      domain_name      = "test2.com"
      validation_scope = "HOST"
    },
    {
      domain_name      = "test3.com"
      validation_scope = "HOST"
    },
    {
      domain_name      = "test4.com"
      validation_scope = "HOST"
    },
  ]
}
