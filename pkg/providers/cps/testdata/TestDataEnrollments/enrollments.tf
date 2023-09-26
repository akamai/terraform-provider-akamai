provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cps_enrollments" "test" {
  contract_id = "testing"
}
