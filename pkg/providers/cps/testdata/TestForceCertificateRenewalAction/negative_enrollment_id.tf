provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "terraform_data" "trigger" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.akamai_cps_force_certificate_renewal.test]
    }
  }
}

action "akamai_cps_force_certificate_renewal" "test" {
  config {
    enrollment_id = -1
  }
}
