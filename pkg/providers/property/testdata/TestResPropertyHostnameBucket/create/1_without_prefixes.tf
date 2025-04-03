provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_hostname_bucket" "test" {
  property_id = "111"
  contract_id = "222"
  group_id    = "333"
  network     = "STAGING"
  hostnames = {
    "www.test.hostname.0.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "444"
    },
  }
}
