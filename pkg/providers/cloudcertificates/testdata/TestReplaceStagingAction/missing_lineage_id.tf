provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "trigger" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.akamai_cloudcertificates_replace_staging.test]
    }
  }
}

action "akamai_cloudcertificates_replace_staging" "test" {
  config {
    generation_id = 3100
  }
}
