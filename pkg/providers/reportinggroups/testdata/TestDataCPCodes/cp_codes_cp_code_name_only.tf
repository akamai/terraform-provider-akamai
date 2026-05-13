provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_cp_codes" "test" {
  cp_code_name = "Other CP Code"
}

