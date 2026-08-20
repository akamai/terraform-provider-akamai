provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "s" {
  log_type = "CDN"
  active   = false

  stream_name = "test-cdn-stream-with-service-ids"

  dataset_fields = [2000]
  properties     = [1]
  service_ids    = [101]

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
