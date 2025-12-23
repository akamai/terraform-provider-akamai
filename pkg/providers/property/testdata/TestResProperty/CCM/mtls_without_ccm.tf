provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  contract_id = "ctr_C-0N7RAC7"
  group_id    = "grp_12345"
  product_id  = "prd_Object_Delivery"

  hostnames {
    cname_from             = "example.com"
    cname_to               = "example.com.edgekey.net"
    cert_provisioning_type = "DEFAULT"
    mtls {
      ca_set_id          = "524125"
      check_client_ocsp  = true
      send_ca_set_client = true
    }
  }
}
