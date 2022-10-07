provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cps_upload_certificate" "test" {
  enrollment_id                          = 2
  certificate_rsa_pem                    = "-----BEGIN CERTIFICATE REQUEST UPDATED-----\n...\n-----END CERTIFICATE REQUEST UPDATED-----"
  acknowledge_post_verification_warnings = false
  acknowledge_change_management          = false
}
