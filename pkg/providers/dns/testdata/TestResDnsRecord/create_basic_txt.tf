provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_dns_record" "txt_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "TXT"
  active     = true
  ttl        = 300
  target     = ["Hel\\lo\"world", "\"extralongtargetwhichis\" \"intwoseparateparts\""]
}

