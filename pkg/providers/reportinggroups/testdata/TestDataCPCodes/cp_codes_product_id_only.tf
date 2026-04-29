provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_cp_codes" "test" {
  product_id = "Site_Accel::Site_Accel"
}

