provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_domain" "testdomain" {
  name                      = "gtm_terra_testdomain.akadns.net"
  type                      = "weighted"
  contract                  = "1-2ABCDEF"
  comment                   = "Edit Property test_property"
  group                     = "123ABC"
  load_imbalance_percentage = 10.0
  email_notification_list   = ["email1@nomail.com", "email2@nomail.com", "email3@nomail.com"]
}
