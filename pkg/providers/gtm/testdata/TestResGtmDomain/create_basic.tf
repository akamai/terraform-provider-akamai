provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_gtm_domain" "testdomain" {
  name                      = "gtm_terra_testdomain.akadns.net"
  type                      = "weighted"
  contract                  = "1-2ABCDEF"
  comment                   = "Test"
  group                     = "123ABC"
  load_imbalance_percentage = 10.0
}
