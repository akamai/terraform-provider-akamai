provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

locals {
  entries  = [for i in range(0, 999) : "www.test.hostname.${i}.com.edgesuite.net"]
  entries2 = [for i in range(1000, 1098) : "www.test.hostname.${i}.com.edgesuite.net"]
}

resource "akamai_property_hostname_bucket" "test" {
  property_id = "prp_111"
  contract_id = "ctr_222"
  group_id    = "grp_333"
  network     = "STAGING"
  hostnames = {
    for entry in concat(local.entries, local.entries2) :
    entry => {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_444"
    }
  }
}
