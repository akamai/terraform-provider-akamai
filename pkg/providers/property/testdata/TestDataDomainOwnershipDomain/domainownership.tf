provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_domainownership_domain" "testdomain" {
  domain_name      = "example.com"
  validation_scope = "DOMAIN"
}