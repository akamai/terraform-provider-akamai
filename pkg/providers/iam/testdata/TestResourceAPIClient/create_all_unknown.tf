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
  credential          = {}
  group_access = random_integer.api_access.result == 1 ? {
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
    } : {
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
  ip_acl = random_integer.api_access.result == 1 ? {
    enable = true
    cidr   = ["128.5.6.5/24"]
    } : {
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
  purge_options = random_integer.api_access.result == 1 ? {
    can_purge_by_cp_code   = true
    can_purge_by_cache_tag = true
    cp_code_access = {
      all_current_and_new_cp_codes = false
      cp_codes                     = [101]
    }
    } : {
    can_purge_by_cp_code   = true
    can_purge_by_cache_tag = true
    cp_code_access = {
      all_current_and_new_cp_codes = false
      cp_codes                     = [101]
    }
  }
}

resource "random_integer" "api_access" {
  min = 0
  max = 1
}