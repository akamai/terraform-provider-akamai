provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_notification_prods" "test" {}
