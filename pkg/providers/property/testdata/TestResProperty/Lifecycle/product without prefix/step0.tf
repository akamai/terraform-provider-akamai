provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  contract_id = "ctr_0"
  group_id    = "grp_0"
  product     = "0"

  hostnames = {

    "cnameTo": "to.test.domain",
    "cnameFrom": "from.test.domain",
    "certProvisioningType": "CPS_MANAGED"
  }
}
