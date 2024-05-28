provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_attack_payload_logging" "test" {
  config_id              = 43253
  attack_payload_logging = file("testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingDisabled.json")
}