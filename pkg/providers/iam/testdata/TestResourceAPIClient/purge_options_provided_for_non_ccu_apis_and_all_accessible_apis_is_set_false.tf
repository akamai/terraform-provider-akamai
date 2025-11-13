provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_api_client" "test" {
  authorized_users    = ["mw+2"]
  client_type         = "CLIENT"
  client_name         = "mw+2_1"
  notification_emails = ["mw+2@example.com"]
  client_description  = "Test API Client"
  credential          = {}
  api_access = {
    all_accessible_apis = false
    apis = [
      {
        api_id       = 5583
        access_level = "READ-ONLY"
      },
      {
        api_id       = 5802
        access_level = "READ-WRITE"
      }
    ]
  }
  group_access = {
    clone_authorized_user_groups = true
  }
  purge_options = {
    can_purge_by_cp_code   = false
    can_purge_by_cache_tag = false
    cp_code_access = {
      all_current_and_new_cp_codes = false
      cp_codes                     = [101]
    }
  }
}