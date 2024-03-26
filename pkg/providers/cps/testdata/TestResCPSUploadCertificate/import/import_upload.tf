provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_upload_certificate" "import" {
  enrollment_id                          = 1
  certificate_rsa_pem                    = "-----BEGIN CERTIFICATE ECDSA -----\n...\n-----END CERTIFICATE ECDSA -----"
  trust_chain_rsa_pem                    = "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA -----\n...\n-----END CERTIFICATE TRUST-CHAIN ECDSA -----"
  acknowledge_post_verification_warnings = false
  auto_approve_warnings                  = []
  acknowledge_change_management          = false
}
