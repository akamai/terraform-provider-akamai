provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_hostname_coverage" "test" {

}

