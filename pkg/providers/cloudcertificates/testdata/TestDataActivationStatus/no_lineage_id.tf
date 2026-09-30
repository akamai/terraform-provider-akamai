provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudcertificates_activation_status" "test" {
  activation_id = 100
}
