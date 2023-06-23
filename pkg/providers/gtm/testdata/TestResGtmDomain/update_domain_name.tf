provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_domain" "testdomain" {
  name                      = "gtm_terra_testdomain.akadns.net-updated"
  type                      = "weighted"
  contract                  = "1-2ABCDEF"
  comment                   = "Test"
  group                     = "123ABC"
  load_imbalance_percentage = 20.0
}
