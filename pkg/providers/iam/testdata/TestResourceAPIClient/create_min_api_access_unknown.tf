provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_api_client" "test" {
  authorized_users = ["mw+2"]
  client_type      = "CLIENT"
  client_name      = "mw+2_1"
  lock             = true
  credential       = {}
  group_access = {
    clone_authorized_user_groups = false
    groups = [
      {
        group_id = 123
        role_id  = 340
      }
    ]
  }
  api_access = random_integer.api_access.result == 1 ? {
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
    } : {
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
}

resource "random_integer" "api_access" {
  min = 0
  max = 1
}