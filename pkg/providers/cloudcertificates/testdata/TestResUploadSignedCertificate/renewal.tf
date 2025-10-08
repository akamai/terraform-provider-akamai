provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudcertificates_upload_signed_certificate" "upload" {
  certificate_id         = "23456"
  signed_certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
testrenewedsignedcertificate
-----END CERTIFICATE-----
EOT
}
