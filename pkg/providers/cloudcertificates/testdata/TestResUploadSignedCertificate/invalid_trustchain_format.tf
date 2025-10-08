provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload_signed_certificate" "upload" {
  certificate_id         = "12345"
  signed_certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
testsignedcertificate
-----END CERTIFICATE-----
EOT
  trust_chain_pem        = <<EOT
-----FOO BAR-----
testtrustchaincertificate1
-----END CERTIFICATE-----
EOT
}
