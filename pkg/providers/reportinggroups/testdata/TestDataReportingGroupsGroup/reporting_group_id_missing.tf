provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_group" "test" {
}

