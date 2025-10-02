provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_validation" "test" {
  domains = concat([
    for i in range(1, 1000) : {
      domain_name = "test${i}.example.com"
    }
    ], [
    {
      domain_name      = "test1001.example.com"
      validation_scope = "HOST"
    }
  ])
}
