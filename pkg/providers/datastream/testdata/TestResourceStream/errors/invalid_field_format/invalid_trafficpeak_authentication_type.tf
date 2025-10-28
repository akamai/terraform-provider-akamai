provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "s" {
  active = false
  delivery_configuration {
    format = "JSON"
    frequency {
      interval_in_secs = 30
    }
  }

  contract_id    = "test_contract"
  dataset_fields = [1001]

  group_id    = 1337
  properties  = [1]
  stream_name = "test_stream"

  trafficpeak_connector {
    display_name        = "display_name"
    authentication_type = "NONE"
    compress_logs       = true
    content_type        = "application/json"
    custom_header_name  = "custom_header_name"
    custom_header_value = "custom_header_value"
    endpoint            = "https://demo.trafficpeak.live/ingest/event?table=ABC&token=123"
    user_name           = "user_name"
    password            = "password"
  }
}
