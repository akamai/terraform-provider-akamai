provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "renamed_property"
  contract_id = "ctr_0"
  group_id    = "grp_0"
  product_id  = "prd_0"

  hostnames {
    cname_to = "to2.test.domain"
    cname_from = "from.test.domain"
    cert_provisioning_type = "DEFAULT"
  }
}
