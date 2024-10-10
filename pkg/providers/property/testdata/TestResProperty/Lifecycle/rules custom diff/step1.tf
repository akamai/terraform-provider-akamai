provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  contract_id = "ctr_1"
  group_id    = "grp_2"
  product_id  = "prd_3"

  rules = data.akamai_property_rules_template.rules.json

  hostnames {
    cname_to               = "to.test.domain"
    cname_from             = "from.test.domain"
    cert_provisioning_type = "DEFAULT"
  }
}

data "akamai_property_rules_template" "rules" {
  template_file = "testdata/TestResProperty/Lifecycle/rules custom diff/property-snippets/rules1.json"
}
