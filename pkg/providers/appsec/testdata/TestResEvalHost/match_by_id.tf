provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_eval_hostnames" "test" {
 config_id = 43253
    version = 7
  hostnames = ["example.com"]
}



