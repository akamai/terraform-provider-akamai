provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  group_id    = "grp_0"
  contract_id = "ctr_0"
  product_id = "prd_0"

  hostnames = [{
    cnameTo: "to.test.domain",
    cnameFrom: "from.test.domain",
    certProvisioningType: "CPS_MANAGED"
  }]
}
