provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  contract_id = "ctr_0"
  group_id    = "grp_0"
  product_id  = "0"

  hostnames = {
    "from.test.domain" = "to.test.domain"
  }
}
