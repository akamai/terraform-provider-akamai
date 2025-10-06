provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_domainownership_domains" "test" {
}