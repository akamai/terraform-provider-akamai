provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_datastream" "s" {
  active = false
  config {
    format = "JSON"
    frequency {
      time_in_sec = 30
    }
  }

  contract_id = "test_contract"
  dataset_fields_ids = [
    1001
  ]
  group_id = 1337
  property_ids = [
    1,
  ]
  stream_name   = "test_stream"
  stream_type   = "RAW_LOGS"
  template_name = "EDGE_LOGS"

  gcs_connector {
    bucket               = "bucket"
    connector_name       = "connector_name"
    path                 = "path"
    private_key          = "private_key"
    project_id           = "project_id"
    service_account_name = "service_account_name"
  }
}
