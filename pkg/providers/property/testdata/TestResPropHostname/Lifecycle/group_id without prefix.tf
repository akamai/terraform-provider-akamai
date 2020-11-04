provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_hostnames" "test" {
  property_id = "prp_0"
  contract_id = "ctr_0"
  group_id    = "0"

  names = {
    "test.domain" = "ehn_0"
  }
}
