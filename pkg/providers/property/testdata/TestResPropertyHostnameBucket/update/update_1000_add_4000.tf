provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

locals {
  entries  = [for i in range(0, 1000) : "www.test.hostname.${i}.com.edgesuite.net"]
  entries2 = [for i in range(1000, 2000) : "www.test.hostname.${i}.com.edgesuite.net"]
  entries3 = [for i in range(2000, 3000) : "www.test.hostname.${i}.com.edgesuite.net"]
  entries4 = [for i in range(3000, 4000) : "www.test.hostname.${i}.com.edgesuite.net"]
  entries5 = [for i in range(4000, 5000) : "www.test.hostname.${i}.com.edgesuite.net"]
}

resource "akamai_property_hostname_bucket" "test" {
  property_id = "prp_111"
  contract_id = "ctr_222"
  group_id    = "grp_333"
  network     = "STAGING"
  hostnames = {
    for entry in concat(local.entries, local.entries2, local.entries3, local.entries4, local.entries5) :
    entry => {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_555"
    }
  }
}
