provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "splunk_stream" {
  active = false
  delivery_configuration {
    format = "JSON"
    frequency {
      interval_in_secs = 30
    }
  }

  contract_id    = "test_contract"
  dataset_fields = [1001, 1002]

  group_id    = 1337
  properties  = [1]
  stream_name = "test_stream"

  splunk_connector {
    compress_logs         = false
    display_name          = "splunk_test_connector_name"
    event_collector_token = "splunk_event_collector_token"
    endpoint              = "splunk_url"
  }
}
