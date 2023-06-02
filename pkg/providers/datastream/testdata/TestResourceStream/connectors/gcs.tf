provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_datastream" "s" {
  active = false
  delivery_configuration {
    format = "JSON"
    frequency {
      interval_in_secs = 30
    }
  }

  contract_id = "test_contract"
  dataset_fields = [
    1001
  ]
  group_id = 1337
  properties = [
    1,
  ]
  stream_name = "test_stream"

  gcs_connector {
    bucket               = "bucket"
    display_name         = "connector_name"
    path                 = "path"
    private_key          = "private_key"
    project_id           = "project_id"
    service_account_name = "service_account_name"
  }
}
