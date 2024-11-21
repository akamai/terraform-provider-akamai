provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_zone" "secondary_test_zone" {
  contract       = "ctr1"
  zone           = "secondaryexampleterraform.io"
  type           = "secondary"
  comment        = "This is an updated test secondary zone"
  sign_and_serve = false
  group          = "grp1"
  masters        = ["1.1.1.1"]
  outbound_zone_transfer {
    acl            = ["192.0.2.156/24"]
    enabled        = true
    notify_targets = ["192.0.2.192"]
    tsig_key {
      algorithm = "hmac-sha1"
      name      = "other.com.akamai.com"
      secret    = "fakeSecretajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
    }
  }
  tsig_key {
    algorithm = "hmac-sha512"
    name      = "other.com.akamai.com"
    secret    = "fakeSecretjVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
  }
}