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

  s3_connector {
    access_key        = "s3_test_access_key"
    bucket            = "s3_test_bucket"
    display_name      = "s3_test_connector_name"
    path              = "s3_test_path"
    region            = "s3_test_region"
    secret_access_key = "s3_test_secret_key"
  }
}
