provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_iam_notification_prods" "test" {}
