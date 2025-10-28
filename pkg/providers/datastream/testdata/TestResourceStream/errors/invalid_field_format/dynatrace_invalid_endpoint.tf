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

  dynatrace_connector {
    display_name        = "display_name"
    endpoint            = "endpoint"
    api_token           = "api_token"
    custom_header_name  = "custom_header_name"
    custom_header_value = "custom_header_value"
  }
}
