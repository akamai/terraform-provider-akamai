provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property" "test" {
  name                = "test_property"
  contract_id         = "1"
  group_id            = "grp_2"
  product_id          = "prd_3"
  use_hostname_bucket = "true"

  hostnames {
    cname_to               = "to.test.domain"
    cname_from             = "from.test.domain"
    cert_provisioning_type = "DEFAULT"
  }
}
