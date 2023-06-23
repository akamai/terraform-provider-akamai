provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cps_deployments" "test" {
  enrollment_id = 123
}