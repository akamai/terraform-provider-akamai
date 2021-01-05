provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_api_hostname_coverage" "hostname_coverage" {
}

