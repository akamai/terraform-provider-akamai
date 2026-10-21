provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "trigger" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.akamai_cloudcertificates_renew.test]
    }
  }
}

action "akamai_cloudcertificates_renew" "test" {
  config {
    lineage_id = -1
  }
}
