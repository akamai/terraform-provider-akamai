provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  contract_id = "ctr_0"
  group_id    = "grp_0"
  product  = "prd_0"

  rules = data.akamai_property_rules_template.akarules.json

  hostnames {
    cname_to= "to.test.domain"
    cname_from="from.test.domain"
    cert_provisioning_type= "DEFAULT"
  }

}

data "akamai_property_rules_template" "akarules" {
  template_file = "testdata/TestResProperty/Lifecycle/property-snippets/rules1.json"
}
