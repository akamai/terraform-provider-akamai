provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_api_client" "test" {
  authorized_users = ["mw+2"]
  client_type      = "CLIENT"
  client_name      = "mw+2_1"
  lock             = false
  credential       = {}
  group_access = {
    clone_authorized_user_groups = true
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
}