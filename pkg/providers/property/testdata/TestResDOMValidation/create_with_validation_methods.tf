provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_validation" "test" {
  domains = [
    {
      domain_name       = "test1.example.com"
      validation_scope  = "HOST"
      validation_method = "DNS_CNAME"
    },
    {
      domain_name       = "test2.example.com"
      validation_scope  = "DOMAIN"
      validation_method = "DNS_TXT"
    },
    {
      domain_name       = "test3.example.com"
      validation_scope  = "WILDCARD"
      validation_method = "HTTP"
    },
  ]
}
