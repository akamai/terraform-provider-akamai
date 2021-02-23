provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  contract_id = "ctr_0"
  group       = "grp_0"
  product_id  = "prd_0"

  hostnames = [
    {
    cnameTo: "to.test.domain",
    cnameFrom: "from.test.domain",
    certProvisioningType: "CPS_MANAGED"
  }
]

}
