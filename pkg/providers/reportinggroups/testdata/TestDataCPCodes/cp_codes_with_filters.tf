provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_reportinggroups_cp_codes" "test" {
  contract_id  = "1-2ABCDE"
  group_id     = "grp_67890"
  product_id   = "Web_Exp::Ion_Na"
  cp_code_name = "My CP Code"
}

