provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_clientlist_activation" "activation_ASN_LIST_1" {
  version                 = 2
  network                 = "STAGING"
  comments                = "Activation Comments"
  notification_recipients = ["user@example.com"]
  siebel_ticket_id        = "ABC-12345"
}
