provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_url_evasion_defense" "test" {
  config_id    = 43253
  status       = "disabled"
  bypass_lists = ["123_TEST", "456_TEST"]

  rules = [
    {
      rule_id            = 950001
      action             = "deny"
      condition_operator = "OR"
      conditions = [
        {
          type           = "extensionMatch"
          extensions     = ["php", "jsp"]
          positive_match = true
        },
        {
          type           = "filenameMatch"
          filenames      = ["index.php"]
          positive_match = true
        },
        {
          type           = "hostMatch"
          hosts          = ["example.com"]
          positive_match = true
        },
        {
          type           = "pathMatch"
          paths          = ["/admin/*"]
          positive_match = true
        },
        {
          type           = "requestMethodMatch"
          methods        = ["GET", "POST"]
          positive_match = true
        },
        {
          type           = "ipMatch"
          ips            = ["1.2.3.4"]
          use_headers    = true
          positive_match = true
        },
        {
          type           = "clientListMatch"
          client_lists   = ["123_TEST"]
          use_headers    = true
          positive_match = true
        },
        {
          type                 = "requestHeaderMatch"
          header               = "X-Test"
          value                = "abc*"
          value_case_sensitive = true
          value_wildcard       = true
          positive_match       = true
        },
        {
          type                 = "uriQueryMatch"
          name                 = "param"
          value                = "val*"
          name_case_sensitive  = true
          value_case_sensitive = true
          value_wildcard       = true
          positive_match       = true
        }
      ]
    }
  ]
}

