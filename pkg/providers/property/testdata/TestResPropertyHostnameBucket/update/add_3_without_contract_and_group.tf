provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_hostname_bucket" "test" {
  property_id = "prp_111"
  network     = "STAGING"
  hostnames = {
    "www.test.hostname.0.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_444"
    },
    "www.test.hostname.1.com.edgesuite.net" : {
      cert_provisioning_type = "DEFAULT"
      edge_hostname_id       = "ehn_555"
    },
    "www.test.hostname.2.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_555"
    },
    "www.test.hostname.3.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_555"
    },
  }
}
