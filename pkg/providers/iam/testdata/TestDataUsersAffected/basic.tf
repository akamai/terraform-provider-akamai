provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_users_affected_by_moving_group" "test" {
  source_group_id      = 123
  destination_group_id = 321
}