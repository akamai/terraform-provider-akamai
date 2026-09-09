provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_clientlist_activation" "activation_ASN_LIST_1" {
  list_id                 = "12_AB"
  version                 = 3
  network                 = "STAGING"
  comments                = "Activation Comments Updated"
  notification_recipients = ["update_user@example.com"]
  siebel_ticket_id        = "UPDATED-12345"
}
