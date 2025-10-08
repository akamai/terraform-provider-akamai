provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_test" "cert_id" {
  input = "12345"
}

resource "akamai_test" "ack" {
  input = "true"
}

resource "akamai_test" "cert_pem" {
  input = <<EOT
-----BEGIN CERTIFICATE-----
testsignedcertificate
-----END CERTIFICATE-----
EOT
}

resource "akamai_test" "chain_pem" {
  input = <<EOT
-----BEGIN CERTIFICATE-----
testtrustchaincertificate1
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
testtrustchaincertificate2
-----END CERTIFICATE-----
EOT
}

resource "akamai_cloudcertificates_upload_signed_certificate" "upload" {
  certificate_id         = akamai_test.cert_id.output
  acknowledge_warnings   = tobool(akamai_test.ack.output)
  signed_certificate_pem = akamai_test.cert_pem.output
  trust_chain_pem        = akamai_test.chain_pem.output
}
