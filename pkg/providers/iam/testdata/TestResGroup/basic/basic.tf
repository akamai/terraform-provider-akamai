provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_iam_group" "test" {
  parent_group_id = 1
  name            = "test"
}