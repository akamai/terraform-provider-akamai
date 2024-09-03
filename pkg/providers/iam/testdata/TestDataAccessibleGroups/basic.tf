provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_accessible_groups" "groups" {
  username = "user1"
}