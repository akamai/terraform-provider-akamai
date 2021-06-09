provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cps_dv_validation" "dv_validation" {
  enrollment_id = 1
}