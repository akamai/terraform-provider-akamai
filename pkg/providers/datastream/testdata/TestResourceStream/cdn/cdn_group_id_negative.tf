provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "s" {
  log_type = "CDN"
  active   = false

  stream_name = "test-cdn-stream-group-id-negative"
  group_id    = "-1"

  delivery_configuration {
    field_delimiter = "SPACE"
    format          = "STRUCTURED"
    frequency {
      interval_in_secs = 30
    }
    upload_file_prefix = "pre"
    upload_file_suffix = "suf"
  }

  dataset_fields      = [1001, 1002, 2000, 2001]
  notification_emails = ["test_email1@akamai.com"]
  properties          = [1, 2, 3]

  s3_connector {
    access_key        = "s3_test_access_key"
    bucket            = "s3_test_bucket"
    display_name      = "s3_test_connector_name"
    path              = "s3_test_path"
    region            = "s3_test_region"
    secret_access_key = "s3_test_secret_key"
  }
}
