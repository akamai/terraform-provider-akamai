provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cps_upload_certificate" "test" {
  enrollment_id                          = 2
  certificate_rsa_pem                    = "-----BEGIN CERTIFICATE RSA REQUEST-----\n...\n-----END CERTIFICATE RSA REQUEST-----"
  trust_chain_rsa_pem                    = "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST-----\n...\n-----END CERTIFICATE TRUST-CHAIN RSA REQUEST-----"
  acknowledge_post_verification_warnings = false
  auto_approve_warnings = [
    "CERTIFICATE_ADDED_TO_TRUST_CHAIN",
    "CERTIFICATE_ALREADY_LOADED",
    "CERTIFICATE_DATA_BLANK_OR_MISSING",
  ]
  acknowledge_change_management = true
}
