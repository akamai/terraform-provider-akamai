provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_clientside_protection_configuration" "test" {
  config_id          = 43253
  security_policy_id = "test_policy_id"

  edge_injection = {
    load_script_async = true
  }

  edge_test_parameters = {
    disable_injection_key   = "disable-key"
    disable_injection_value = "disable-value"
    force_injection_key     = "force-key"
    force_injection_value   = "force-value"
  }

  injection_criteria = {
    injection_path_includes = "/include-path"
    injection_path_excludes = "/exclude-path"
    injection_percent       = 100
    injection_policy        = "head"
  }
}