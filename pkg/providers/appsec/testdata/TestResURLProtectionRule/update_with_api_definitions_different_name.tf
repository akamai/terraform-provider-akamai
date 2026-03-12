provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}
resource "akamai_appsec_url_protection_rule" "test" {
  config_id          = 43007
  name               = "API Protection Rule test"
  description        = "Updated API Protection"
  max_rate_threshold = 195

  api_definitions = [{
    api_definition_id   = 3216157
    defined_resources   = true
    resource_ids        = []
    undefined_resources = true
    }
  ]
}
