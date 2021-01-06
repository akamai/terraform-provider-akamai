provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_eval_hostnames" "eval_hostnames" {
  config_id = 43253
    version = 7
}

