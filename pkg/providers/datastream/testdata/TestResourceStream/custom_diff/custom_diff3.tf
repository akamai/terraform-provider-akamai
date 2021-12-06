provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_datastream" "s" {
  active = true
  config {
    delimiter = "SPACE"
    format    = "STRUCTURED"
    frequency {
      time_in_sec = 30
    }
    upload_file_suffix = "suf"
  }

  contract_id = "test_contract"
  dataset_fields_ids = [
    1000,
    1001,
    1002
  ]
  email_ids = [
    "test_email1@akamai.com",
    "test_email2@akamai.com",
  ]
  group_id = 1337
  property_ids = [
    1,
    2,
    3
  ]
  stream_name   = "test_stream"
  stream_type   = "RAW_LOGS"
  template_name = "EDGE_LOGS"

  sumologic_connector {
    collector_code = "sumologic_collector_code"
    connector_name = "sumologic_connector_name"
    endpoint       = "sumologic_endpoint"
  }
}
