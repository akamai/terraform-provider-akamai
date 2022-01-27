provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_cloudlets_application_load_balancer" "alb" {
  origin_id      = "test_origin"
  description    = "test description updated"
  balancing_type = "PERFORMANCE"
  data_centers {
    cloud_server_host_header_override = false
    cloud_service                     = true
    country                           = "US"
    continent                         = "NA"
    latitude                          = 102.78108
    longitude                         = -116.07064
    percent                           = 100
    liveness_hosts = [
    "tf.test"]
    hostname          = "test-hostname"
    state_or_province = "MA"
    city              = "Boston"
    origin_id         = "test_origin"
  }
  liveness_settings {
    port        = 1234
    protocol    = "HTTP"
    path        = "/status"
    host_header = "header"
    additional_headers = {
      abc = "123"
    }
    interval        = 10
    request_string  = "test_request_string"
    response_string = "test_response_string"
    timeout         = 60
  }
}