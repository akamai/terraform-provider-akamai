provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudwrapper_location" "test" {
  location_name = "US West"
  traffic_type  = "WEB_ENHANCED_TLS"
}