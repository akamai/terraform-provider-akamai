provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_ip_allowlist" "test" {
  enable = true
}
