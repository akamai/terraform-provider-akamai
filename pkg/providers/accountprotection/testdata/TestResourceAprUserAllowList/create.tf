provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_apr_user_allow_list" "test" {
  config_id = 43253
  user_allow_list = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}