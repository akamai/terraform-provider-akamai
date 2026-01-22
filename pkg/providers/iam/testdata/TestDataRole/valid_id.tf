provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}


data "akamai_iam_role" "test" {
  role_id = 12345
}
