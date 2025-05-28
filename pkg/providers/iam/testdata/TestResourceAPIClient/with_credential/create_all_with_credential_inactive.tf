provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_api_client" "test" {
  authorized_users    = ["mw+2"]
  client_type         = "CLIENT"
  client_name         = "mw+2_1"
  notification_emails = ["mw+2@example.com"]
  client_description  = "Test API Client"
  lock                = false
  credential = {
    description = "Test API Client Credential"
    expires_on  = "2025-06-13T14:48:07Z"
    status      = "INACTIVE"
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
  ip_acl = {
    enable = true
    cidr   = ["128.5.6.5/24"]
  }
  api_access = {
    all_accessible_apis = false
    apis = [
      {
        api_id       = 5580
        access_level = "READ-ONLY"
      },
      {
        api_id       = 5801
        access_level = "READ-WRITE"
      }
    ]
  }
  purge_options = {
    can_purge_by_cp_code   = true
    can_purge_by_cache_tag = true
    cp_code_access = {
      all_current_and_new_cp_codes = false
      cp_codes                     = [101]
    }
  }
}