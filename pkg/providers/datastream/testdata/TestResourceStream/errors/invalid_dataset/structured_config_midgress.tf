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
  dataset_fields = [2051, 1001, 1002]

  group_id    = 1337
  properties  = [1]
  stream_name = "test_stream"

  s3_connector {
    bucket            = "s3_bucket"
    display_name      = "s3_display_name"
    path              = "s3_path"
    access_key        = "s3_access_key"
    region            = "s3_region"
    secret_access_key = "s3_secret_key"
  }
}
