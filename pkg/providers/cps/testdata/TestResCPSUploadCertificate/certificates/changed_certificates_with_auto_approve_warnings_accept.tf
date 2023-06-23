provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_upload_certificate" "test" {
  enrollment_id                          = 2
  certificate_rsa_pem                    = "-----BEGIN CERTIFICATE RSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE RSA REQUEST UPDATED-----"
  certificate_ecdsa_pem                  = "-----BEGIN CERTIFICATE ECDSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE ECDSA REQUEST UPDATED-----"
  trust_chain_rsa_pem                    = "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE TRUST-CHAIN RSA REQUEST UPDATED-----"
  trust_chain_ecdsa_pem                  = "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE TRUST-CHAIN ECDSA REQUEST UPDATED-----"
  acknowledge_post_verification_warnings = false
  auto_approve_warnings = [
    "CERTIFICATE_ADDED_TO_TRUST_CHAIN",
    "CERTIFICATE_ALREADY_LOADED",
    "CERTIFICATE_DATA_BLANK_OR_MISSING",
    "CERTIFICATE_HAS_NULL_ISSUER",
  ]
  acknowledge_change_management = true
}