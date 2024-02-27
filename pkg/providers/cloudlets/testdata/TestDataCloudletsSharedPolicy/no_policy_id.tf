provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_shared_policy" "test" {}