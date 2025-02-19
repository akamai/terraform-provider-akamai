provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostname_activation" "activation" {
  property_id            = "prp_1"
  contract_id            = "ctr_1"
  group_id               = "grp_1"
  hostname_activation_id = "atv_1"
  include_hostnames      = "true"
}