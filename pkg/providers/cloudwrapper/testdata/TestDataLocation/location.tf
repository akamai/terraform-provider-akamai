provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudwrapper_location" "test" {
  location_name = "US West"
  traffic_type  = "LIVE_VOD"
}