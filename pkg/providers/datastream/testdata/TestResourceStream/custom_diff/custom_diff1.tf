provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "s" {
  active = true
  delivery_configuration {
    field_delimiter = "SPACE"
    format          = "STRUCTURED"
    frequency {
      interval_in_secs = 30
    }
    upload_file_prefix = "pre"
    upload_file_suffix = "suf"
  }

  contract_id = "test_contract"
  dataset_fields = [
    1000,
    1001,
    1002
  ]
  notification_emails = [
    "test_email1@akamai.com",
    "test_email2@akamai.com",
  ]
  group_id = 1337
  properties = [
    1,
    2,
    3
  ]
  stream_name = "test_stream"

  sumologic_connector {
    collector_code = "sumologic_collector_code"
    display_name   = "sumologic_connector_name"
    endpoint       = "sumologic_endpoint"
  }
}
