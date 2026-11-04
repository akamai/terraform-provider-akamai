provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "trigger" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.akamai_cloudcertificates_rollback.test]
    }
  }
}

action "akamai_cloudcertificates_rollback" "test" {
  config {
    lineage_id    = 800002
    generation_id = 3001

    timeouts = {
      invoke = "1h"
    }
  }
}
