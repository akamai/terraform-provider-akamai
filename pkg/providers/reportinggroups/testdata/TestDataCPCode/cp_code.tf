provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_cp_code" "test" {
  cp_code_id = 12345
}

