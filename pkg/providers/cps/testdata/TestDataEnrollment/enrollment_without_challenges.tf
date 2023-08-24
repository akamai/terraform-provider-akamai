provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cps_enrollment" "test" {
  enrollment_id = 1
}
