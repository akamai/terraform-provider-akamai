provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

# intelligent_load_shedding is set with hits_per_sec and categories ONLY.
# custom_criteria is intentionally omitted (null in plan).
# Since custom_criteria is no longer Computed, the provider must NOT
# populate it from the API response when the plan didn't include it.
resource "akamai_appsec_url_protection_policy" "test" {
  config_id          = 43007
  name               = "URL Protection"
  description        = "URL Protection"
  max_rate_threshold = 195

  hostname_paths = [{
    hostname = "custom.com"
    paths    = ["/asd", "/my-test-path"]
  }]

  intelligent_load_shedding = {
    hits_per_sec = 150
    categories   = ["BOTS", "CLOUD_PROVIDERS", "PROXIES", "TOR_EXIT_NODES", "PLATFORM_DDOS_INTELLIGENCE"]
    # custom_criteria deliberately omitted
  }
}

