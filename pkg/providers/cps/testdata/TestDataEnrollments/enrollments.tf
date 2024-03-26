provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cps_enrollments" "test" {
  contract_id = "testing"
}
