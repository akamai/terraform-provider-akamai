provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cp_code" "akacpcode" {
  contract_id = var.contractid
  group_id    = var.groupid
  product_id  = var.product
  name        = var.cp_code
}

variable "groupid" {
  type    = string
  default = "grp_22"
}

variable "contractid" {
  type    = string
  default = "ctr_11"
}

variable "product" {
  type    = string
  default = "prd_3"
}

variable "cp_code" {
  type    = string
  default = "test-ft-cp-code"
}