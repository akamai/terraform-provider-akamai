provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_zone" "multi_signer_invalid_test_zone" {
  contract       = "ctr1"
  zone           = "multisignerinvalidexampleterraform.io"
  type           = "primary"
  comment        = "This is an invalid test zone with multi-signer DNSSEC but no sign_and_serve"
  sign_and_serve = false
  group          = "grp1"

  multi_provider_dnssec {
    enabled = true
  }
}
