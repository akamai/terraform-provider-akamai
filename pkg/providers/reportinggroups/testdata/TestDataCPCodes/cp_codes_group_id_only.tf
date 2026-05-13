provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_cp_codes" "test" {
  group_id = "grp_67890"
}

