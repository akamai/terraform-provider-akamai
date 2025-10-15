// filepath: /Users/nisb/g/projects/tfm/terraform-provider-akamai/pkg/providers/appsec/testdata/TestResActivations/match_by_id_multiple_emails.tf
provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_activations" "test" {
  config_id           = 43253
  version             = 7
  network             = "STAGING"
  note                = "Test Notes"
  notification_emails = ["user1@example.com", "user2@example.com"]
}
