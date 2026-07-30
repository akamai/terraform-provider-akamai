provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "s" {
  log_type = "APPSEC"
  active   = false

  stream_name = "test-app-sec-stream-whitespace-contract"
  contract_id = "   "
  group_id    = "42"

  notification_emails = [
    "nobody@akamai.com"
  ]

  app_sec_configs = [16536]
  properties      = []

  trafficpeak_connector {
    authentication_type = "BASIC"
    display_name        = "TrafficPeakTest"
    endpoint            = "https://example.com/ingest/event?table=unit_test&token=1234"
    content_type        = "application/json"
    compress_logs       = true
    user_name           = "username"
    password            = "password"
  }

  delivery_configuration {
    format = "JSON"

    frequency {
      interval_in_secs = 30
    }
  }
}
