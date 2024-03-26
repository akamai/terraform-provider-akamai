provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_group" "test" {
  parent_group_id = 7
  name            = "another test"
}