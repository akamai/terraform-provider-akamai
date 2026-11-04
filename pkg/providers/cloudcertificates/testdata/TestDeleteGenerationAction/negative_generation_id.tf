provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "trigger" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.akamai_cloudcertificates_delete_generation.test]
    }
  }
}

action "akamai_cloudcertificates_delete_generation" "test" {
  config {
    lineage_id    = 920001
    generation_id = -1
  }
}
