provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_datastream" "seed" {
  log_type = "CDN"
  active   = false

  stream_name    = "seed-cdn-stream"
  dataset_fields = [2000]
  properties     = [1]

  trafficpeak_connector {
    authentication_type = "BASIC"
    display_name        = "TrafficPeakSeed"
    endpoint            = "https://example.com/ingest/event?table=seed&token=1234"
    content_type        = "application/json"
    compress_logs       = true
    user_name           = "username"
    password            = "password"
  }

  delivery_configuration {
    format = "JSON"

    frequency {
      interval_in_secs = 30
    }
  }
}

data "akamai_datastream_answerx_service_ids" "unknown" {
  # unknown during plan, so this data source is deferred and service_ids remain unknown
  contract_id = akamai_datastream.seed.id
}

resource "akamai_datastream" "s" {
  log_type = "APPSEC"
  active   = false

  stream_name     = "appsec-unknown-forbidden-service-ids"
  group_id        = "42"
  contract_id     = "test_contract"
  app_sec_configs = [16536]
  service_ids     = [for s in data.akamai_datastream_answerx_service_ids.unknown.service_ids : s.id]

  trafficpeak_connector {
    authentication_type = "BASIC"
    display_name        = "TrafficPeakTarget"
    endpoint            = "https://example.com/ingest/event?table=target&token=1234"
    content_type        = "application/json"
    compress_logs       = true
    user_name           = "username"
    password            = "password"
  }

  delivery_configuration {
    format = "JSON"

    frequency {
      interval_in_secs = 30
    }
  }
}
