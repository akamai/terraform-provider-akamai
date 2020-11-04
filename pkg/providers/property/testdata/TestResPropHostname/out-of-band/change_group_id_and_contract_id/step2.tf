provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_hostnames" "test" {
  property_id = "prp_0"
  contract_id = "ctr_1"
  group_id    = "grp_1"

  names = {
    "test.domain" = "ehn_0"
  }
}
