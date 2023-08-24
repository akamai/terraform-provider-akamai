provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_iam_states" "test" {
  country = "test country"
}
