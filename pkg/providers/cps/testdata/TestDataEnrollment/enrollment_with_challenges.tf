provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cps_enrollment" "test" {
  enrollment_id = 2
}