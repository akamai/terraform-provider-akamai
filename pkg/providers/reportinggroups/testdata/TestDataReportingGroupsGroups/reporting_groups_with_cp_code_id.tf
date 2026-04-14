provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_groups" "test" {
  cp_code_id = "111"
}

