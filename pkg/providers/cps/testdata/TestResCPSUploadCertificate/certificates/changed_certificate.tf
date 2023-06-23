provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_upload_certificate" "test" {
  enrollment_id                          = 2
  certificate_rsa_pem                    = "-----BEGIN CERTIFICATE RSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE RSA REQUEST UPDATED-----"
  acknowledge_post_verification_warnings = true
  acknowledge_change_management          = true
}
