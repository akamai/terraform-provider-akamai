provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_hostname_bucket" "test" {
  property_id = "prp_111"
  contract_id = "ctr_222"
  group_id    = "grp_333"
  network     = "STAGING"
  hostnames = {
    "www.test.hostname.0.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_444"
    },
  }
}
