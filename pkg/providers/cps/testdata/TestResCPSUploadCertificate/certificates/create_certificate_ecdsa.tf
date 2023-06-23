provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_upload_certificate" "test" {
  enrollment_id                          = 2
  certificate_ecdsa_pem                  = "-----BEGIN CERTIFICATE ECDSA REQUEST-----\n...\n-----END CERTIFICATE ECDSA REQUEST-----"
  acknowledge_post_verification_warnings = true
  acknowledge_change_management          = false
}
