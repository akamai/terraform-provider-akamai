provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_states" "test" {
  country = "test country"
}
