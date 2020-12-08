provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  group_id    = "grp_0"
  contract_id = "ctr_0"
  product_id = "prd_0"

  hostnames = {
    "from.test.domain" = "to.test.domain"
  }
}
