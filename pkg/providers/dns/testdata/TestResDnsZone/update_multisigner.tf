provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_dns_zone" "multi_signer_test_zone" {
  contract              = "ctr1"
  zone                  = "multisignerexampleterraform.io"
  type                  = "primary"
  comment               = "This is an updated test zone with multi-signer DNSSEC"
  sign_and_serve        = true
  group                 = "grp1"
  multi_provider_dnssec = true
}
