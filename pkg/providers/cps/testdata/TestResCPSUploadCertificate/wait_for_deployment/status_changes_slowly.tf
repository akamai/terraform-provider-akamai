provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cps_upload_certificate" "test" {
  enrollment_id                          = 2
  certificate_rsa_pem                    = "-----BEGIN CERTIFICATE RSA REQUEST-----\n...\n-----END CERTIFICATE RSA REQUEST-----"
  acknowledge_post_verification_warnings = true
  wait_for_deployment                    = true
}
