provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_bot_detection" "test" {
  detection_name = "Test name 3"
}