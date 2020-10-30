provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_hostnames" "test" {
  contract_id = "ctr_0"
  group_id    = "grp_0"

  names = {
    "test.domain" = "ehn_0"
  }
}
