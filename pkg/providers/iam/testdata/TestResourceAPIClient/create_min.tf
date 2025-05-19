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
  api_access = {
    all_accessible_apis = false
    apis = [
      {
        api_id            = 5580
        api_name          = "Search Data Feed"
        description       = "Search Data Feed"
        endpoint          = "/search-portal-data-feed-api/"
        documentation_url = "/"
        access_level      = "READ-ONLY"
      },
      {
        api_id            = 5801
        api_name          = "EdgeWorkers"
        description       = "EdgeWorkers"
        endpoint          = "/edgeworkers/"
        documentation_url = "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html"
        access_level      = "READ-WRITE"
      }
    ]
  }
}