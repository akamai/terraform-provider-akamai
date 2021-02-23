provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "renamed property"
  contract_id = "ctr_0"
  group_id    = "grp_0"
  product_id  = "prd_0"

  hostnames =  [{
    cnameTo: "to2.test.domain",
    cnameFrom: "from.test.domain",
    certProvisioningType: "CPS_MANAGED"
  }]
}
