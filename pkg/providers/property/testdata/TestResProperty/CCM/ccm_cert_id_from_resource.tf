provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "lineage_id" {
  input = "164877"
}

resource "akamai_property" "test" {
  name        = "test_property"
  contract_id = "ctr_C-0N7RAC7"
  group_id    = "grp_12345"
  product_id  = "prd_Object_Delivery"

  hostnames {
    cname_from             = "example.com"
    cname_to               = "example.com.edgekey.net"
    cert_provisioning_type = "CCM"
    ccm_cert_id            = terraform_data.lineage_id.output
  }
}
