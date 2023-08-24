provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cps_deployments" "test" {
  enrollment_id = 123
}