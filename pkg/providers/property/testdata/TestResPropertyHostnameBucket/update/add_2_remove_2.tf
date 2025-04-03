provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_hostname_bucket" "test" {
  property_id = "prp_111"
  contract_id = "ctr_222"
  group_id    = "grp_333"
  network     = "STAGING"
  hostnames = {
    "www.test.hostname.5.com.edgesuite.net" : {
      cert_provisioning_type = "DEFAULT"
      edge_hostname_id       = "ehn_555"
    },
    "www.test.hostname.0.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_444"
    },
    #    "www.test.hostname.1.com.edgesuite.net": {
    #      cert_provisioning_type = "DEFAULT"
    #      edge_hostname_id = "ehn_444"
    #    },
    "www.test.hostname.2.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_444"
    },
    #    "www.test.hostname.3.com.edgesuite.net": {
    #      cert_provisioning_type = "CPS_MANAGED"
    #      edge_hostname_id = "ehn_555"
    #    },
    "www.test.hostname.6.com.edgesuite.net" : {
      cert_provisioning_type = "DEFAULT"
      edge_hostname_id       = "ehn_666"
    },
    "www.test.hostname.4.com.edgesuite.net" : {
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_444"
    },
  }
}
