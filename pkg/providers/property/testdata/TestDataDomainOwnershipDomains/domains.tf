provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_domainownership_domains" "test" {
}