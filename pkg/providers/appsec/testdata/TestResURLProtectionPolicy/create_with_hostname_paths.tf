provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_url_protection_policy" "test" {
  config_id          = 43007
  name               = "URL Protection"
  description        = "URL Protection"
  max_rate_threshold = 195

  hostname_paths = [{
    hostname = "custom.com"
    paths    = ["/asd", "/my-test-path"]
    }
  ]

  bypass_conditions = [{
    type                 = "RequestHeaderCondition"
    names                = ["my-custom-header"]
    name_wildcard        = false
    values               = ["my-custom-value"]
    value_case_sensitive = false
    value_wildcard       = false
    },
    {
      type   = "NetworkListCondition"
      names  = []
      values = ["12345_10CLIENTLIST", "54321_123"]
    }
  ]
  intelligent_load_shedding = {
    hits_per_sec = 150
    categories   = ["BOTS", "CLOUD_PROVIDERS", "PROXIES", "TOR_EXIT_NODES", "PLATFORM_DDOS_INTELLIGENCE"]
    custom_criteria = [
      {
        type           = "CLIENT_LIST"
        list_ids       = ["12345_10CLIENTLIST", "54321_123"]
        positive_match = true
      }
    ]
  }
}

