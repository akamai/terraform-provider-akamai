provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "answerx_stream" {
  log_type = "ANSWERX"
  active   = false

  contract_id = "test_contract"
  group_id    = 1337
  stream_name = "test_answerx_stream"

  service_ids    = [2320, 510, 2925]
  dataset_fields = [2000]

  trafficpeak_connector {
    authentication_type = "BASIC"
    display_name        = "tp_connector"
    endpoint            = "https://example.com/ingest/event?table=test&token=tok"
    content_type        = "application/json"
    compress_logs       = true
    user_name           = "user"
    password            = "pass"
  }

  delivery_configuration {
    format = "JSON"
    frequency {
      interval_in_secs = 30
    }
  }
}
