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
    tls_configuration {
      cipher_profile              = "ak-akamai-2020q1"
      disallowed_tls_versions     = ["TLSv1_1", "TLSv1"]
      staple_server_ocsp_response = true
      fips_mode                   = true
    }
  }
}
