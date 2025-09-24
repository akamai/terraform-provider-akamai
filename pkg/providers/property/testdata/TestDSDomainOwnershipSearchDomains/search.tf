provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_domainownership_search_domains" "test" {
  domains = [
    {
      domain_name      = "dom1.test"
      validation_scope = "HOST"
    },
    {
      domain_name      = "dom2.test"
      validation_scope = "HOST"
    },
    {
      domain_name      = "dom3.test"
      validation_scope = "HOST"
    },
    {
      domain_name      = "dom4.test"
      validation_scope = "HOST"
    }
  ]
}
