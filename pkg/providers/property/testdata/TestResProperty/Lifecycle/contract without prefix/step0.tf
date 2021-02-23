provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  contract    = "0"
  group_id    = "grp_0"
  product_id  = "prd_0"

  hostnames = [
    {
    cnameTo: "to.test.domain",
    cnameFrom: "from.test.domain",
    certProvisioningType: "CPS_MANAGED"
  }
]

}
