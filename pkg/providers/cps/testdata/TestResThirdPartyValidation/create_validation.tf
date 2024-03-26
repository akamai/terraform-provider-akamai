provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cps_third_party_validation" "third_party_validation" {
  enrollment_id = 1
  sans = [
    "san.test.akamai.com",
  ]
}