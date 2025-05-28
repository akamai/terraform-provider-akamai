provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_api_client" "test" {
  authorized_users    = ["mw+2"]
  client_type         = "CLIENT"
  notification_emails = ["mw+2@example.com"]
  client_name         = "mw+2_1"
  client_description  = "Test API Client"
  api_access = {
    all_accessible_apis = true
  }
  group_access = {
    clone_authorized_user_groups = false
    groups = [
      {
        group_id = 123
        role_id  = 340
        sub_groups = [
          {
            group_id = 333
            role_id  = 540
          }
        ]
      }
    ]
  }
  allow_account_switch = false
  lock                 = false
  ip_acl = {
    enable = false
    cidr   = ["128.5.6.5/24"]
  }
  purge_options = {
    can_purge_by_cp_code   = false
    can_purge_by_cache_tag = false
    cp_code_access = {
      all_current_and_new_cp_codes = false
      cp_codes                     = [101]
    }
  }
  credential = {}
}