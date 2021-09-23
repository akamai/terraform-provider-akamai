provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_application_load_balancer" test {
  origin_id = "application_load_balancer"
}