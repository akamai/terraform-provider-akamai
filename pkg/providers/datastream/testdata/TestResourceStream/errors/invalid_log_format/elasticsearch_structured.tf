provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "splunk_stream" {
  active = false
  delivery_configuration {
    format          = "STRUCTURED"
    field_delimiter = "SPACE"
    frequency {
      interval_in_secs = 30
    }
  }

  contract_id    = "test_contract"
  dataset_fields = [1001, 1002]

  group_id    = 1337
  properties  = [1]
  stream_name = "test_stream"

  elasticsearch_connector {
    display_name = "display_name"
    index_name   = "index_name"
    endpoint     = "endpoint"
    user_name    = "user_name"
    password     = "password"
  }
}
