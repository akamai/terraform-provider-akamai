provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  contract_id = "ctr_0"
  group_id    = "grp_0"
  product_id  = "prd_0"

  hostnames = {
    "from.test.domain" = "to2.test.domain"
  }

  rules = "{\"rules\":{\"name\":\"default\",\"options\":{}}}"
}
