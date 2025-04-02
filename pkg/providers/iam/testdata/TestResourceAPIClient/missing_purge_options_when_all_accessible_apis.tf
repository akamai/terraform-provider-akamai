provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_api_client" "test" {
  authorized_users    = ["mw+2"]
  client_type         = "CLIENT"
  client_name         = "mw+2_1"
  notification_emails = ["mw+2@example.com"]
  client_description  = "Test API Client"
  api_access = {
    all_accessible_apis = true
  }
  group_access = {
    clone_authorized_user_groups = true
  }
}