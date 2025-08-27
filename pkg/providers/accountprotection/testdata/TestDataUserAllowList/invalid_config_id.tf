provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_apr_user_allow_list" "test" {
  config_id = "badconfigid"
}
