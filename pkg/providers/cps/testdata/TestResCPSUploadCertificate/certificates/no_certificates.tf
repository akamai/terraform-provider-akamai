provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cps_upload_certificate" "test" {
  enrollment_id                          = 2
  acknowledge_post_verification_warnings = true
  acknowledge_change_management          = true
  wait_for_deployment                    = true
}
