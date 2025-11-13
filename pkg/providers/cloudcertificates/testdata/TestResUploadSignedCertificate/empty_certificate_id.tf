provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload_signed_certificate" "upload" {
  certificate_id         = ""
  signed_certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
testsignedcertificate
-----END CERTIFICATE-----
EOT
}
